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
	mw "sdmht/sdmht/api/http/middleware"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"
)

const (
	ConnServeAddrKey = "sdmht_conn_addr"
)

func getConnServeAddr(ctx context.Context) string {
	if v := ctx.Value(kitx.MetadataKey(ConnServeAddrKey)); v != nil {
		if value, ok := v.([]string); ok {
			return value[0]
		}
	}
	return ""
}

var _ itfs.SignalingService = (*signalingSvc)(nil)

type signalingSvc struct {
	// repo       itfs.EventParticipantRepo
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
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("NewLineup: get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		req.AccountID = operater
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
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("FindLineup: get operator fail when update event")
			return nil, lib.NewError(lib.ErrInvalidArgument, "get account id fail")
		}
		query.FilterByAccountID = operater
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
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("UpdateLineup: get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		req.AccountID = operater
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
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		req.AccountID = operater
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

	player := &entity.Player{}
	player.ID = req.AccountID
	player.Scene = entity.NewScene(lineup.CardLibrarys, req.Positions)
	player.Units = units

	match := &entity.Match{}
	match.ID, _ = s.idGenerator.NextID()
	match.Players = append(match.Players, player)

	s.matchRepo.SetByAccount(player.ID, match.ID)
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

	player := &entity.Player{}
	player.ID = req.AccountID
	player.Scene = entity.NewScene(lineup.CardLibrarys, req.Positions)
	player.Units = units

	match, err := s.matchRepo.Get(req.MatchID)
	if err != nil {
		log.S().Errorw("JoinMatch: get match fail", "err", err)
		return nil, err
	}
	if len(match.Players) >= 2 {
		log.S().Errorw("JoinMatch: get match fail", "err", err)
		return nil, lib.NewError(lib.ErrInternal, "match is already full")
	}

	match.WhoseTurn = int32(rand.Perm(2)[0])
	match.Players = append(match.Players, player)

	s.matchRepo.SetByAccount(player.ID, match.ID)
	err = s.matchRepo.Join(&match)
	if err != nil {
		log.S().Errorw("JoinMatch: join match fail", "err", err)
		return nil, err
	}

	data, err := json.Marshal(&match)
	if err != nil {
		log.S().Errorw("JoinMatch: marshal match fail", "err", err)
		return nil, err
	}
	log.S().Errorw("JoinEvent:MarshalNotice", "data", json.RawMessage(data))
	cli, e := s.connManger.GetConnCli(context.TODO(), req.AccountID)
	if e != nil {
		log.S().Errorw("JoinEvent:Dispatch", "accountid", req.AccountID, "noClient", cli == nil, "err", e)
	}
	go func(match entity.Match, data []byte) {
		accountID := match.Players[0].ID
		_, err := cli.DispatchEventToClient(context.TODO(), accountID, entity.ClientEvent{
			AccountID: accountID,
			Type:      entity.MsgTypeSyncMatch,
			AtTime:    time.Now(),
			Content:   data,
		})
		if err != nil {
			log.S().Errorw("JoinEvent", "err", err)
		}
	}(match, data)

	return &entity.JoinMatchRes{
		Match: match,
	}, nil
}

func (s *signalingSvc) SyncOperate(ctx context.Context, req *entity.SyncOperate) (*entity.SyncOperateRes, error) {
	match, err := s.matchRepo.Get(req.MatchID)
	if err != nil {
		log.S().Errorw("SyncOperate: get match fail", "matchID", req.MatchID, "err", err)
		return nil, err
	}

	syncs := &entity.SyncOperateRes{
		Operates: []*entity.SyncOperate{},
	}

	switch req.Event {
	case entity.OpEventAttack:
		if match.WhoseTurn != req.Operator {
			return nil, lib.NewError(lib.ErrInvalidArgument, "不是你的回合")
		}
		other := match.GetOtherPlayer()
		// 费用检查
		if match.Players[req.Operator].Scene.Cost < entity.DefaultAttachCost {
			return nil, lib.NewError(lib.ErrInvalidArgument, "费用不足")
		}
		// 攻击的迷雾位置
		if match.Players[other].Scene.Squares[req.To] == 0 {
			match.Players[other].Scene.Squares[req.To] += entity.SquareExposeTime
		}
		// 检查迷雾位置的单位
		toUnit := match.Players[other].Scene.UnitsLocation[req.To]
		if toUnit != 0 {
			return syncs, nil
		}
		// 攻击发起者
		// fromUnit := match.Players[req.Operator].Units[req.From]

		for _, unit := range match.Players[other].Units {
			if unit.ID != toUnit {
				continue
			}
		}

	case entity.OpEventMove:
	case entity.OpEventCard:
	case entity.OpEventSkill:
	case entity.OpEventEndRound:
	default:
		log.S().Errorw("SyncOperator", "req", req)
		return nil, lib.NewError(lib.ErrInvalidArgument, "非法操作")
	}

	return syncs, nil
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
