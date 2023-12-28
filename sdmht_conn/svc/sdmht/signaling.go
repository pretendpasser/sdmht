package sdmht

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	account_entity "sdmht/account/svc/entity"
	account "sdmht/account/svc/interfaces"
	"sdmht/lib"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/seq"
	"sdmht/sdmht_conn/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"
)

type ContextKey string

const (
	ConnServeAddrKey ContextKey = "sdmht_conn_addr"
)

func getConnServeAddr(ctx context.Context) string {
	if v := ctx.Value(ConnServeAddrKey); v != nil {
		if value, ok := v.(string); ok {
			return value
		}
	}
	return ""
}

var _ itfs.SignalingService = (*signalingSvc)(nil)

type signalingSvc struct {
	idGenerator seq.IDGenerator
	lineupRepo  itfs.LineupRepo
	unitRepo    itfs.UnitRepo
	matchRepo   itfs.MatchRepo
	accountSvc  account.Service
	connManger  *ConnManager
}

func NewSignalingService(idGenerator seq.IDGenerator,
	lineupRepo itfs.LineupRepo,
	unitRepo itfs.UnitRepo,
	matchRepo itfs.MatchRepo,
	accountSvc account.Service,
	connManger *ConnManager) itfs.SignalingService {
	return &signalingSvc{
		idGenerator: idGenerator,
		lineupRepo:  lineupRepo,
		unitRepo:    unitRepo,
		matchRepo:   matchRepo,
		accountSvc:  accountSvc,
		connManger:  connManger,
	}
}

func (s *signalingSvc) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	err := s.accountSvc.Register(ctx, &account_entity.RegisterReq{
		WechatID: req.WeChatID,
		UserName: req.UserName,
	})
	if err != nil {
		if !errors.Is(err, lib.NewError(lib.ErrInvalidArgument, "account exist")) {
			log.S().Errorw("Login: account register fail", "err", err)
			return nil, err
		}
	}

	res, err := s.accountSvc.Login(ctx, &account_entity.LoginReq{
		WechatID: req.WeChatID,
	})
	if err != nil {
		log.S().Errorw("Login: account login fail", "err", err)
		return nil, err
	}

	serveAddr := getConnServeAddr(ctx)
	err = s.connManger.User2ConnRepo().Add(ctx, res.Account.ID, serveAddr)
	if err != nil {
		log.S().Errorw("Login: add user client fail", "err", err)
		return nil, err
	}

	return &entity.LoginRes{
		AccountID: res.Account.ID,
	}, nil
}

func (s *signalingSvc) NewLineup(ctx context.Context, req *entity.NewLineupReq) error {
	if len(req.Units) > 3 || len(req.CardLibrarys) > entity.MaxCardLibrary {
		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
	}

	if req.AccountID == 0 {
		log.S().Error("NewLineup: get account id fail")
		return lib.NewError(lib.ErrInternal, "get account id fail")
	}

	req.Enabled = true
	if len(req.Units) < 3 || len(req.CardLibrarys) < entity.MaxCardLibrary {
		req.Enabled = false
	}
	if req.Name == "" {
		req.Name = "自定义卡组"
	}

	err := s.lineupRepo.Create(ctx, &req.Lineup)
	if err != nil {
		log.S().Errorw("NewLineup: create lineup fail", "err", err)
		return err
	}
	return nil
}

func (s *signalingSvc) FindLineup(ctx context.Context, req *entity.FindLineupReq) (*entity.FindLineupRes, error) {
	query := &entity.LineupQuery{FilterByAccountID: req.AccountID}
	if req.AccountID == 0 {
		log.S().Error("NewLineup: get account id fail")
		return nil, lib.NewError(lib.ErrInternal, "get account id fail")
	}

	total, lineups, err := s.lineupRepo.Find(ctx, query)
	if err != nil {
		log.S().Errorw("FindLineup: find lineup fail", "err", err)
		return nil, err
	}

	return &entity.FindLineupRes{
		Total:   total,
		Lineups: lineups,
	}, nil
}

func (s *signalingSvc) UpdateLineup(ctx context.Context, req *entity.UpdateLineupReq) error {
	if len(req.Units) > 3 || len(req.CardLibrarys) > entity.MaxCardLibrary {
		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
	}

	if req.AccountID == 0 {
		log.S().Error("NewLineup: get account id fail")
		return lib.NewError(lib.ErrInternal, "get account id fail")
	}

	req.Enabled = true
	if len(req.Units) < 3 || len(req.CardLibrarys) < entity.MaxCardLibrary {
		req.Enabled = false
	}
	if req.Name == "" {
		req.Name = "自定义卡组"
	}

	err := s.lineupRepo.Update(ctx, &req.Lineup)
	if err != nil {
		log.S().Errorw("UpdateLineup: update lineup fail", "err", err)
		return err
	}
	return nil
}

func (s *signalingSvc) DeleteLineup(ctx context.Context, req *entity.DeleteLineupReq) error {
	if req.AccountID == 0 {
		log.S().Error("NewLineup: get account id fail")
		return lib.NewError(lib.ErrInternal, "get account id fail")
	}

	err := s.lineupRepo.Delete(ctx, req.AccountID, req.ID)
	if err != nil {
		log.S().Errorw("DeleteLineup: delete lineup fail", "err", err)
		return err
	}
	return nil
}

func (s *signalingSvc) NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error) {
	lineup, err := s.lineupRepo.Get(ctx, req.AccountID, req.LineupID)
	if err != nil {
		log.S().Errorw("NewMatch:get lineup fail",
			"accountid", req.AccountID, "lineupid", req.LineupID,
			"err", err)
		return nil, err
	}
	if !lineup.Enabled {
		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}
	unitsLocation := make(map[int64]int32)
	for _, unitID := range lineup.Units {
		unitsLocation[unitID] = -1
	}

	unitIDs := []int64{}
	for position, unitID := range req.Positions {
		if unitID != 0 && unitsLocation[unitID] == -1 {
			unitsLocation[unitID] = int32(position)
			unitIDs = append(unitIDs, unitID)
		}
	}
	if len(unitIDs) != entity.MaxBaseUnitNum {
		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}

	units, err := s.unitRepo.Get(ctx, unitIDs)
	if err != nil {
		log.S().Errorw("NewMatch: get units fail", "err", err)
		return nil, err
	}

	scene := entity.NewScene(req.AccountID, units, lineup.CardLibrarys, req.Positions)
	match := &entity.Match{}
	match.ID, _ = s.idGenerator.NextID()
	match.Scenes = append(match.Scenes, scene)

	s.matchRepo.SetByAccount(scene.PlayerID, match.ID)
	err = s.matchRepo.New(match)
	if err != nil {
		log.S().Errorw("NewMatch: new match fail", "err", err)
		return nil, err
	}

	return &entity.NewMatchRes{
		MatchID: match.ID,
	}, nil
}

func (s *signalingSvc) GetMatch(ctx context.Context, req *entity.GetMatchReq) (*entity.GetMatchRes, error) {
	matchID := s.matchRepo.GetByAccount(req.AccountID)
	match, err := s.matchRepo.Get(matchID)
	if err != nil {
		log.S().Errorw("GetMatch: get match fail", "err", err)
		return nil, err
	}

	return &entity.GetMatchRes{
		Match: match,
	}, nil
}

func (s *signalingSvc) JoinMatch(ctx context.Context, req *entity.JoinMatchReq) (*entity.JoinMatchRes, error) {
	lineup, err := s.lineupRepo.Get(ctx, req.AccountID, req.LineupID)
	if err != nil {
		log.S().Errorw("JoinMatch: NewMatch get lineup fail",
			"accountid", req.AccountID, "lineupid", req.LineupID,
			"err", err)
		return nil, err
	}
	if !lineup.Enabled {
		log.S().Errorw("JoinMatch: lineup is disabled", "err", err)
		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}
	unitsLocation := make(map[int64]int32)
	for _, unitID := range lineup.Units {
		unitsLocation[unitID] = -1
	}

	unitIDs := []int64{}
	for position, unitID := range req.Positions {
		if unitID != 0 && unitsLocation[unitID] == -1 {
			unitsLocation[unitID] = int32(position)
			unitIDs = append(unitIDs, unitID)
		}
	}
	if len(unitIDs) != entity.MaxBaseUnitNum {
		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}

	units, err := s.unitRepo.Get(ctx, unitIDs)
	if err != nil {
		log.S().Errorw("JoinMatch: NewMatch get units fail", "err", err)
		return nil, err
	}

	scene := entity.NewScene(req.AccountID, units, lineup.CardLibrarys, req.Positions)
	match, err := s.matchRepo.Get(req.MatchID)
	if err != nil {
		log.S().Errorw("JoinMatch: get match fail", "err", err)
		return nil, err
	}
	if len(match.Scenes) >= 2 {
		log.S().Errorw("JoinMatch: get match fail", "err", err)
		return nil, lib.NewError(lib.ErrInternal, "match is already full")
	}
	match.Scenes = append(match.Scenes, scene)
	match.WhoseTurn = int32(rand.Perm(2)[0])
	match.Scenes[match.WhoseTurn].Cost = entity.FirstCost
	match.Scenes[match.WhoseTurn].HandCards = append(match.Scenes[match.WhoseTurn].HandCards, 0)

	data, err := json.Marshal(&match)
	if err != nil {
		log.S().Errorw("JoinMatch: marshal match fail", "err", err)
		return nil, err
	}
	log.S().Infow("JoinEvent:MarshalNotice", "data", json.RawMessage(data))
	cli, e := s.connManger.GetConnCli(context.TODO(), req.AccountID)
	if e != nil {
		log.S().Errorw("JoinEvent:Dispatch", "accountid", req.AccountID, "noClient", cli == nil, "err", e)
	}
	_, err = cli.DispatchEventToClient(context.TODO(), match.Scenes[0].PlayerID, entity.ClientEvent{
		AccountID: match.Scenes[0].PlayerID,
		Type:      entity.MsgTypeSyncMatch,
		AtTime:    time.Now(),
		Content:   data,
	})
	if err != nil {
		log.S().Errorw("JoinEvent: dispatch event fail", "err", err)
		return nil, lib.NewError(lib.ErrInternal, "join event fail")
	}

	s.matchRepo.SetByAccount(scene.PlayerID, match.ID)
	err = s.matchRepo.Join(&match)
	if err != nil {
		log.S().Errorw("JoinMatch: join match fail", "err", err)
		return nil, err
	}

	return &entity.JoinMatchRes{
		Match: match,
	}, nil
}

// 最后一个return之前禁止修改数值
func (s *signalingSvc) SyncOperate(ctx context.Context, req *entity.SyncOperateReq) (*entity.SyncOperateRes, error) {
	slog := log.L().With(kitx.TraceIDField(ctx)).Sugar()
	match, err := s.matchRepo.Get(req.MatchID)
	if err != nil {
		slog.Errorw("SyncOperate: get match fail", "matchID", req.MatchID, "err", err)
		return nil, err
	}
	if match.WhoseTurn != req.Operator {
		slog.Errorw("SyncOperate: not your turn", "operator", req.Operator, "req", req)
		return nil, lib.NewError(lib.ErrInvalidArgument, "不是你的回合")
	}
	otherPlayer := match.GetOtherPlayer()
	otherPlayerID := match.Scenes[otherPlayer].PlayerID

	switch req.Event {
	case entity.OpEventAttack:
		if req.To > 15 || req.To < 0 {
			slog.Errorw("SyncOperate: invalid To", "event", req.Event,
				"to", req.To)
			return nil, lib.NewError(lib.ErrInvalidArgument, "非法目标")
		}
		// 费用检查
		if match.Scenes[req.Operator].Cost < entity.DefaultAttachCost {
			slog.Errorw("SyncOperate: no enough cost", "event", req.Event,
				"cost", match.Scenes[req.Operator].Cost)
			return nil, lib.NewError(lib.ErrInvalidArgument, "费用不足")
		}
		// 攻击发起者
		fromUnit := match.Scenes[req.Operator].Units[req.From]
		if fromUnit == nil {
			slog.Errorw("SyncOperate: fromUnit not exist", "event", req.Event,
				"from", req.From)
			return nil, lib.NewError(lib.ErrInvalidArgument, "单位不存在")
		}
		// 检查迷雾位置的单位
		toUnitID := match.Scenes[otherPlayer].UnitsLocation[req.To]
		if toUnitID != 0 {
			toUnit := match.Scenes[otherPlayer].Units[toUnitID]
			if toUnit == nil {
				slog.Error("SyncOperate: toUnit not exist", "event", req.Event,
					"to", req.To)
				return nil, lib.NewError(lib.ErrInvalidArgument, "目标单位不存在")
			}
			if toUnit.AttackPrevent {
				// 圣盾防止
				match.Scenes[otherPlayer].Units[toUnitID].AttackPrevent = false
			} else {
				attack := fromUnit.Attack
				// 护盾结算
				if toUnit.Defend > 0 && attack > 0 {
					if toUnit.Defend >= attack {
						match.Scenes[otherPlayer].Units[toUnitID].Defend -= attack
						attack = 0
					} else {
						attack -= toUnit.Defend
						match.Scenes[otherPlayer].Units[toUnitID].Defend = 0
					}
				}
				// 生命值结算
				if toUnit.Health > 0 && attack > 0 {
					if toUnit.Health > attack {
						match.Scenes[otherPlayer].Units[toUnitID].Health -= attack
					} else {
						delete(match.Scenes[otherPlayer].Units, toUnitID)
						match.Scenes[otherPlayer].UnitsLocation[req.To] = 0
						if toUnit.Rarity >= 1 && toUnit.Rarity <= 4 {
							match.Scenes[otherPlayer].RetainerAliveNum--
						}
					}
				}
			}
		}
		// 攻击的迷雾位置  暗雾->开雾
		if match.Scenes[otherPlayer].Squares[req.To] == entity.OriginSquare {
			match.Scenes[otherPlayer].Squares[req.To] += entity.SquareExposeTime
		}
		match.Scenes[req.Operator].Cost -= entity.DefaultAttachCost
	case entity.OpEventMove:
		// 检查目的地是否合法
		if req.To > 15 || req.To < 0 || match.Scenes[req.Operator].UnitsLocation[req.To] != 0 {
			slog.Errorw("SyncOperate: invalid From or To", "event", req.Event,
				"to", req.To, "locationTo", match.Scenes[req.Operator].UnitsLocation[req.To])
			return nil, lib.NewError(lib.ErrInvalidArgument, "非法目标")
		}
		// 移动发起者
		fromUnit := match.Scenes[req.Operator].Units[req.From]
		if fromUnit == nil {
			slog.Errorw("SyncOperate: fromUnit not exist", "event", req.Event,
				"from", req.From)
			return nil, lib.NewError(lib.ErrInvalidArgument, "单位不存在")
		}
		// 检查是否可以移动
		if fromUnit.BaseAttribute.BaseNoMove || fromUnit.NoMove > 0 ||
			fromUnit.Move == 0 || fromUnit.IsMoving == -1 {
			slog.Errorw("SyncOperate: not allowed to move", "event", req.Event,
				"BaseNoMove", fromUnit.BaseAttribute.BaseNoMove, "NoMove", fromUnit.NoMove,
				"Move", fromUnit.Move, "IsMoving", fromUnit.IsMoving,
				"Cost", match.Scenes[req.Operator].Cost)
			return nil, lib.NewError(lib.ErrInvalidArgument, "禁止移动")
		}
		if match.Scenes[req.Operator].Cost == 0 && fromUnit.IsMoving == 0 {
			slog.Errorw("SyncOperate: no cost to move", "event", req.Event,
				"Cost", match.Scenes[req.Operator].Cost)
			return nil, lib.NewError(lib.ErrInvalidArgument, "费用不足")
		}
		// 移动发起者的位置
		var fromUnitLocation int64 = 0
		for position, unitID := range match.Scenes[req.Operator].UnitsLocation {
			if unitID == fromUnit.ID {
				if int64(position) == req.To {
					slog.Errorw("SyncOperate: From is equal to To",
						"fromLocation", fromUnitLocation, "to", req.To)
					return nil, lib.NewError(lib.ErrInvalidArgument, "移动始终点相同")
				}
				fromUnitLocation = int64(position)
				break
			}
			if position == len(match.Scenes[req.Operator].UnitsLocation) {
				slog.Errorw("SyncOperate: lose sync between unit and it's location", "event", req.Event,
					"from", req.From, "to", req.To, "unitLocation", match.Scenes[req.Operator].UnitsLocation)
				return nil, lib.NewError(lib.ErrInvalidArgument, "单位与位置不同步")
			}
		}
		err := entity.CheckMoveing(fromUnitLocation, req.To)
		if err != nil {
			slog.Errorw("SyncOperate:check moving fail", "event", req.Event,
				"fromLoc", fromUnitLocation, "toLoc", req.To)
			return nil, err
		}
		if req.From != match.Scenes[req.Operator].LastMoveUnitID {
			lastMoveUnitID := match.Scenes[req.Operator].LastMoveUnitID
			if match.Scenes[req.Operator].Units[lastMoveUnitID] != nil {
				match.Scenes[req.Operator].Units[lastMoveUnitID].IsMoving = -1
			}
		}
		match.Scenes[req.Operator].UnitsLocation[fromUnitLocation] = 0
		match.Scenes[req.Operator].UnitsLocation[req.To] = fromUnit.ID
		match.Scenes[req.Operator].Units[req.From].Move -= 1
		if match.Scenes[req.Operator].Units[req.From].IsMoving == 0 {
			match.Scenes[req.Operator].Cost -= 1
		}
		match.Scenes[req.Operator].Units[req.From].IsMoving = 1
		match.Scenes[req.Operator].LastMoveUnitID = fromUnit.ID
	case entity.OpEventCard:
		length := len(match.Scenes[req.Operator].HandCards)
		for i, card := range match.Scenes[req.Operator].HandCards {
			if card == req.From {
				max := length - 1
				match.Scenes[req.Operator].HandCards[i], match.Scenes[req.Operator].HandCards[max] =
					match.Scenes[req.Operator].HandCards[max], match.Scenes[req.Operator].HandCards[i]
				match.Scenes[req.Operator].HandCards = match.Scenes[req.Operator].HandCards[:max]
				break
			}
			if i == length {
				slog.Errorw("SyncOperate: no exist this card", "event", req.Event,
					"from", req.From, "to", req.To)
				return nil, lib.NewError(lib.ErrInvalidArgument, "卡牌不存在")
			}
		}
	case entity.OpEventSkill:
	case entity.OpEventEndRound:
		// 切换当前回合
		match.WhoseTurn = int32(match.GetOtherPlayer())
		// 对手回合开始
		match.Scenes[otherPlayer].NextRound()
	default:
		slog.Errorw("SyncOperator: no this event", "req", req)
		return nil, lib.NewError(lib.ErrInvalidArgument, "非法操作")
	}

	lastMoveUnitID := match.Scenes[req.Operator].LastMoveUnitID
	if lastMoveUnitID != 0 && req.Event != entity.OpEventMove {
		if match.Scenes[req.Operator].Units[lastMoveUnitID] != nil {
			match.Scenes[req.Operator].Units[lastMoveUnitID].IsMoving = -1
		}
		match.Scenes[req.Operator].LastMoveUnitID = 0
	}

	// check match end
	for i, scene := range match.Scenes {
		masterID := scene.MasterID
		if scene.Units[masterID].Health == 0 {
			match.Winner = match.Scenes[1-i].PlayerID
		}
	}
	s.matchRepo.Set(&match)

	data, err := json.Marshal(&match)
	if err != nil {
		slog.Errorw("SyncOperator: marshal match fail", "err", err)
		return nil, err
	}
	slog.Infow("SyncOperator:MarshalNotice", "data", json.RawMessage(data))
	cli, e := s.connManger.GetConnCli(context.TODO(), otherPlayerID)
	if e != nil {
		slog.Errorw("SyncOperator:Dispatch", "accountid", otherPlayerID, "noClient", cli == nil, "err", e)
	}
	go func(accountID uint64, data []byte) {
		_, err := cli.DispatchEventToClient(context.TODO(), accountID, entity.ClientEvent{
			AccountID: accountID,
			Type:      entity.MsgTypeSyncMatch,
			AtTime:    time.Now(),
			Content:   data,
		})
		if err != nil {
			slog.Errorw("SyncOperator", "err", err)
		}
	}(otherPlayerID, data)

	return &entity.SyncOperateRes{
		Match: match,
	}, nil
}

func (s *signalingSvc) KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error {
	slog := log.L().With(kitx.TraceIDField(ctx)).Sugar()
	slog.Infow("KeepAlive:Req", "params", req)
	// return s.repo.Update(ctx, 0, req.Operator, "status", entity.UserOnline)
	return nil
}

func (s *signalingSvc) Offline(ctx context.Context, req *entity.LogoutReq) error {
	slog := log.L().With(kitx.TraceIDField(ctx)).Sugar()
	slog.Infow("Offline:Req", "params", req)

	// info, err := s.repo.Get(ctx, 0, req.Operator)
	// if err != nil {
	// 	slog.Errorw("Offline:GetJoinedEvent", "params", req, "err", err)
	// 	return err
	// }

	// return s.LeaveEvent(ctx, entity.LeaveEventReq{Operator: req.Operator, EventID: info.EventID})
	return nil
}
