package sdmht

import (
	"context"

	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/lib/seq"
	mw "sdmht/sdmht/api/http/middleware"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.Service = (*service)(nil)

type service struct {
	idGenerator seq.IDGenerator

	lineupRepo itfs.LineupRepo
	unitRepo   itfs.UnitRepo
	matchRepo  itfs.MatchRepo
}

func NewService(idGenerator seq.IDGenerator,
	lineupRepo itfs.LineupRepo,
	unitRepo itfs.UnitRepo,
	matchRepo itfs.MatchRepo) *service {
	return &service{
		idGenerator: idGenerator,
		lineupRepo:  lineupRepo,
		unitRepo:    unitRepo,
		matchRepo:   matchRepo,
	}
}

func (s *service) CreateLineup(ctx context.Context, lineup *entity.Lineup) error {
	if lineup.AccountID == 0 {
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		lineup.AccountID = operater
	}

	if len(lineup.Units) > 3 || len(lineup.CardLibrarys) > entity.MaxCardLibrary {
		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
	}
	if len(lineup.Units) < 3 || len(lineup.CardLibrarys) < entity.MaxCardLibrary {
		lineup.Enabled = false
	}
	if lineup.Name == "" {
		lineup.Name = "自定义卡组"
	}

	err := s.lineupRepo.Create(ctx, lineup)
	if err != nil {
		log.S().Errorw("create lineup fail", "err", err)
		return err
	}

	return nil
}

func (s *service) FindLineup(ctx context.Context, query *entity.LineupQuery) (int, []*entity.Lineup, error) {
	operater, ok := mw.GetAccountIDFromContext(ctx)
	if !ok {
		log.S().Error("get operator fail when update event")
		return 0, nil, lib.NewError(lib.ErrInternal, "get account id fail")
	}
	query.FilterByAccountID = operater

	total, res, err := s.lineupRepo.Find(ctx, query)
	if res != nil {
		log.S().Errorw("find lineup fail", "err", err)
		return 0, nil, err
	}

	return total, res, nil
}

func (s *service) GetLineup(ctx context.Context, req *entity.GetLineupReq) (*entity.Lineup, error) {
	if req.AccountID == 0 {
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return nil, lib.NewError(lib.ErrInternal, "get account id fail")
		}
		req.AccountID = operater
	}

	res, err := s.lineupRepo.Get(ctx, req.ID, req.AccountID)
	if err != nil {
		log.S().Errorw("get lineup fail", "err", err)
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateLineup(ctx context.Context, lineup *entity.Lineup) error {
	if lineup.AccountID == 0 {
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		lineup.AccountID = operater
	}

	if len(lineup.Units) > 3 || len(lineup.CardLibrarys) > entity.MaxCardLibrary {
		return lib.NewError(lib.ErrInvalidArgument, "numbers of units or cards is over max")
	}
	if len(lineup.Units) < 3 || len(lineup.CardLibrarys) < entity.MaxCardLibrary {
		lineup.Enabled = false
	}
	if lineup.Name == "" {
		lineup.Name = "自定义卡组"
	}

	err := s.lineupRepo.Update(ctx, lineup)
	if err != nil {
		log.S().Errorw("update lineup fail", "err", err)
		return err
	}

	return nil
}

func (s *service) DeleteLineup(ctx context.Context, req *entity.DeleteLineupReq) error {
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
		log.S().Errorw("delete lineup fail", "err", err)
		return err
	}

	return nil
}

func (s *service) NewMatch(ctx context.Context, req *entity.NewMatchReq) (uint64, error) {
	lineup, err := s.lineupRepo.Get(ctx, req.AccountID, req.LineupID)
	if err != nil {
		log.S().Errorw("NewMatch get lineup fail", "err", err)
		return 0, err
	}
	if !lineup.Enabled {
		return 0, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}
	unitsLocation := make(map[uint64]int32)
	for _, unitID := range lineup.Units {
		unitsLocation[unitID] = -1
	}

	unitIDs := []uint64{}
	for position, unitID := range req.Positions {
		if unitID != 0 && unitsLocation[unitID] == -1 {
			unitsLocation[unitID] = int32(position)
			unitIDs = append(unitIDs, unitID)
		}
	}
	if len(unitIDs) != entity.MaxBaseUnitNum {
		return 0, lib.NewError(lib.ErrInvalidArgument, "unit number invalid")
	}

	units, err := s.unitRepo.Get(ctx, unitIDs)
	if err != nil {
		log.S().Errorw("NewMatch get units fail", "err", err)
		return 0, err
	}
	for i, unit := range units {
		units[i].Location = unitsLocation[unit.ID]
	}

	player := &entity.Player{}
	player.ID = req.AccountID
	player.Scene = entity.NewScene(lineup.CardLibrarys)
	player.Units = units

	match := &entity.Match{}
	match.ID, _ = s.idGenerator.NextID()
	match.Players = append(match.Players, player)

	err = s.matchRepo.New(match)
	if err != nil {
		log.S().Errorw("NewMatch fail", "err", err)
		return 0, err
	}

	return match.ID, nil
}

func (s *service) JoinMatch() {
	// curRoundID = rand.Perm(2)[0]
}

func (s *service) EndMatch() {}
