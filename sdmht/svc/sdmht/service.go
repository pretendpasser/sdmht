package sdmht

import (
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.Service = (*service)(nil)

type service struct {
	// idGenerator seq.IDGenerator

	// lineupRepo itfs.LineupRepo
	// unitRepo   itfs.UnitRepo
	// matchRepo  itfs.MatchRepo
}

func NewService() itfs.Service { return service{} }

// func NewService(idGenerator seq.IDGenerator,
// 	lineupRepo itfs.LineupRepo,
// 	unitRepo itfs.UnitRepo,
// 	matchRepo itfs.MatchRepo) *service {
// 	return &service{
// 		idGenerator: idGenerator,
// 		lineupRepo:  lineupRepo,
// 		unitRepo:    unitRepo,
// 		matchRepo:   matchRepo,
// 	}
// }

// func (s *service) CreateLineup(ctx context.Context, lineup *entity.Lineup) error {
// 	if lineup.AccountID == 0 {
// 		operater, ok := mw.GetAccountIDFromContext(ctx)
// 		if !ok {
// 			log.S().Error("get operator fail when update event")
// 			return lib.NewError(lib.ErrInternal, "get account id fail")
// 		}
// 		lineup.AccountID = operater
// 	}

// 	lineup.Enabled = true
// 	if len(lineup.Units) > 3 || len(lineup.CardLibrarys) > entity.MaxCardLibrary {
// 		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
// 	}
// 	if len(lineup.Units) < 3 || len(lineup.CardLibrarys) < entity.MaxCardLibrary {
// 		lineup.Enabled = false
// 	}
// 	if lineup.Name == "" {
// 		lineup.Name = "自定义卡组"
// 	}

// 	err := s.lineupRepo.Create(ctx, lineup)
// 	if err != nil {
// 		log.S().Errorw("create lineup fail", "err", err)
// 		return err
// 	}

// 	return nil
// }

// func (s *service) FindLineup(ctx context.Context, query *entity.LineupQuery) (int, []*entity.Lineup, error) {
// 	if query.FilterByAccountID == 0 {
// 		operater, ok := mw.GetAccountIDFromContext(ctx)
// 		if !ok {
// 			log.S().Error("get operator fail when update event")
// 			return 0, nil, lib.NewError(lib.ErrInternal, "get account id fail")
// 		}
// 		query.FilterByAccountID = operater
// 	}

// 	total, res, err := s.lineupRepo.Find(ctx, query)
// 	if err != nil {
// 		log.S().Errorw("find lineup fail", "err", err)
// 		return 0, nil, err
// 	}

// 	return total, res, nil
// }

// func (s *service) GetLineup(ctx context.Context, req *entity.GetLineupReq) (*entity.Lineup, error) {
// 	if req.AccountID == 0 {
// 		operater, ok := mw.GetAccountIDFromContext(ctx)
// 		if !ok {
// 			log.S().Error("get operator fail when update event")
// 			return nil, lib.NewError(lib.ErrInternal, "get account id fail")
// 		}
// 		req.AccountID = operater
// 	}

// 	res, err := s.lineupRepo.Get(ctx, req.ID, req.AccountID)
// 	if err != nil {
// 		log.S().Errorw("get lineup fail", "err", err)
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (s *service) UpdateLineup(ctx context.Context, lineup *entity.Lineup) error {
// 	if lineup.AccountID == 0 {
// 		operater, ok := mw.GetAccountIDFromContext(ctx)
// 		if !ok {
// 			log.S().Error("get operator fail when update event")
// 			return lib.NewError(lib.ErrInternal, "get account id fail")
// 		}
// 		lineup.AccountID = operater
// 	}

// 	lineup.Enabled = true
// 	if len(lineup.Units) > 3 || len(lineup.CardLibrarys) > entity.MaxCardLibrary {
// 		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
// 	}
// 	if len(lineup.Units) < 3 || len(lineup.CardLibrarys) < entity.MaxCardLibrary {
// 		lineup.Enabled = false
// 	}
// 	if lineup.Name == "" {
// 		lineup.Name = "自定义卡组"
// 	}

// 	err := s.lineupRepo.Update(ctx, lineup)
// 	if err != nil {
// 		log.S().Errorw("update lineup fail", "err", err)
// 		return err
// 	}

// 	return nil
// }

// func (s *service) DeleteLineup(ctx context.Context, req *entity.DeleteLineupReq) error {
// 	if req.AccountID == 0 {
// 		operater, ok := mw.GetAccountIDFromContext(ctx)
// 		if !ok {
// 			log.S().Error("get operator fail when update event")
// 			return lib.NewError(lib.ErrInternal, "get account id fail")
// 		}
// 		req.AccountID = operater
// 	}

// 	err := s.lineupRepo.Delete(ctx, req.AccountID, req.ID)
// 	if err != nil {
// 		log.S().Errorw("delete lineup fail", "err", err)
// 		return err
// 	}

// 	return nil
// }

// func (s *service) NewMatch(ctx context.Context, req *entity.NewMatchReq) (uint64, error) {
// 	lineup, err := s.lineupRepo.Get(ctx, req.AccountID, req.LineupID)
// 	if err != nil {
// 		log.S().Errorw("NewMatch get lineup fail",
// 			"accountid", req.AccountID, "lineupid", req.LineupID,
// 			"err", err)
// 		return 0, err
// 	}
// 	if !lineup.Enabled {
// 		return 0, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
// 	}
// 	unitsLocation := make(map[int64]int32)
// 	for _, unitID := range lineup.Units {
// 		unitsLocation[unitID] = -1
// 	}

// 	unitIDs := []int64{}
// 	for position, unitID := range req.Positions {
// 		if unitID != 0 && unitsLocation[unitID] == -1 {
// 			unitsLocation[unitID] = int32(position)
// 			unitIDs = append(unitIDs, unitID)
// 		}
// 	}
// 	if len(unitIDs) != entity.MaxBaseUnitNum {
// 		return 0, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
// 	}

// 	units, err := s.unitRepo.Get(ctx, unitIDs)
// 	if err != nil {
// 		log.S().Errorw("NewMatch get units fail", "err", err)
// 		return 0, err
// 	}

// 	player := &entity.Player{}
// 	player.ID = req.AccountID
// 	player.Scene = entity.NewScene(lineup.CardLibrarys, req.Positions)
// 	player.Units = units

// 	match := &entity.Match{}
// 	match.ID, _ = s.idGenerator.NextID()
// 	match.Players = append(match.Players, player)

// 	s.matchRepo.SetByAccount(player.ID, match.ID)
// 	err = s.matchRepo.New(match)
// 	if err != nil {
// 		log.S().Errorw("NewMatch fail", "err", err)
// 		return 0, err
// 	}

// 	return match.ID, nil
// }

// func (s *service) GetMatch(ctx context.Context, req *entity.GetMatchReq) (*entity.Match, error) {
// 	matchID := s.matchRepo.GetByAccount(req.AccountID)
// 	match, err := s.matchRepo.Get(matchID)
// 	if err != nil {
// 		log.S().Errorw("GetMatch: get match fail", "err", err)
// 		return nil, err
// 	}
// 	return &match, nil
// }

// func (s *service) JoinMatch(ctx context.Context, req *entity.JoinMatchReq) (*entity.Match, error) {
// 	lineup, err := s.lineupRepo.Get(ctx, req.AccountID, req.LineupID)
// 	if err != nil {
// 		log.S().Errorw("NewMatch get lineup fail",
// 			"accountid", req.AccountID, "lineupid", req.LineupID,
// 			"err", err)
// 		return nil, err
// 	}
// 	if !lineup.Enabled {
// 		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
// 	}
// 	unitsLocation := make(map[int64]int32)
// 	for _, unitID := range lineup.Units {
// 		unitsLocation[unitID] = -1
// 	}

// 	unitIDs := []int64{}
// 	for position, unitID := range req.Positions {
// 		if unitID != 0 && unitsLocation[unitID] == -1 {
// 			unitsLocation[unitID] = int32(position)
// 			unitIDs = append(unitIDs, unitID)
// 		}
// 	}
// 	if len(unitIDs) != entity.MaxBaseUnitNum {
// 		return nil, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
// 	}

// 	units, err := s.unitRepo.Get(ctx, unitIDs)
// 	if err != nil {
// 		log.S().Errorw("NewMatch get units fail", "err", err)
// 		return nil, err
// 	}

// 	player := &entity.Player{}
// 	player.ID = req.AccountID
// 	player.Scene = entity.NewScene(lineup.CardLibrarys, req.Positions)
// 	player.Units = units

// 	match, err := s.matchRepo.Get(req.MatchID)
// 	if err != nil {
// 		log.S().Errorw("JoinMatch: get match fail", "err", err)
// 		return nil, err
// 	}
// 	if len(match.Players) >= 2 {
// 		log.S().Errorw("JoinMatch: get match fail", "err", err)
// 		return nil, lib.NewError(lib.ErrInternal, "match is already full")
// 	}

// 	match.WhoseTurn = int32(rand.Perm(2)[0])
// 	match.Players = append(match.Players, player)

// 	s.matchRepo.SetByAccount(player.ID, match.ID)
// 	err = s.matchRepo.Join(&match)
// 	if err != nil {
// 		log.S().Errorw("JoinMatch: join match fail", "err", err)
// 		return nil, err
// 	}

// 	return &match, nil
// }

// func (s *service) SyncOperate(ctx context.Context, req *entity.SyncOperate) (*entity.SyncOperateRes, error) {
// 	match, err := s.matchRepo.Get(req.MatchID)
// 	if err != nil {
// 		log.S().Errorw("SyncOperate: get match fail", "matchID", req.MatchID, "err", err)
// 		return nil, err
// 	}

// 	syncs := &entity.SyncOperateRes{
// 		Operates: []*entity.SyncOperate{},
// 	}

// 	switch req.Event {
// 	case entity.OpEventAttack:
// 		if match.WhoseTurn != req.Operator {
// 			return nil, lib.NewError(lib.ErrInvalidArgument, "不是你的回合")
// 		}
// 		other := match.GetOtherPlayer()
// 		// 费用检查
// 		if match.Players[req.Operator].Scene.Cost < entity.DefaultAttachCost {
// 			return nil, lib.NewError(lib.ErrInvalidArgument, "费用不足")
// 		}
// 		// 攻击的迷雾位置
// 		if match.Players[other].Scene.Squares[req.To] == 0 {
// 			match.Players[other].Scene.Squares[req.To] += entity.SquareExposeTime
// 		}
// 		// 检查迷雾位置的单位
// 		toUnit := match.Players[other].Scene.UnitsLocation[req.To]
// 		if toUnit != 0 {
// 			return syncs, nil
// 		}
// 		// 攻击发起者
// 		// fromUnit := match.Players[req.Operator].Units[req.From]

// 		for _, unit := range match.Players[other].Units {
// 			if unit.ID != toUnit {
// 				continue
// 			}
// 		}

// 	case entity.OpEventMove:
// 	case entity.OpEventCard:
// 	case entity.OpEventSkill:
// 	case entity.OpEventEndRound:
// 	default:
// 		log.S().Errorw("SyncOperator", "req", req)
// 		return nil, lib.NewError(lib.ErrInvalidArgument, "非法操作")
// 	}
// 	return nil, nil
// }
