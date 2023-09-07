package sdmht

import (
	"context"

	"sdmht/lib"
	"sdmht/lib/log"
	mw "sdmht/sdmht/api/http/middleware"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.Service = (*service)(nil)

type service struct {
	lineupRepo itfs.LineupRepo
}

func NewService(lineupRepo itfs.LineupRepo) *service {
	return &service{
		lineupRepo: lineupRepo,
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

func (s *service) GetLineup(ctx context.Context, id uint64, accountID uint64) (*entity.Lineup, error) {
	if accountID == 0 {
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return nil, lib.NewError(lib.ErrInternal, "get account id fail")
		}
		accountID = operater
	}

	res, err := s.lineupRepo.Get(ctx, id, accountID)
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

func (s *service) DeleteLineup(ctx context.Context, id uint64, accountID uint64) error {
	if accountID == 0 {
		operater, ok := mw.GetAccountIDFromContext(ctx)
		if !ok {
			log.S().Error("get operator fail when update event")
			return lib.NewError(lib.ErrInternal, "get account id fail")
		}
		accountID = operater
	}

	err := s.lineupRepo.Delete(ctx, accountID, id)
	if err != nil {
		log.S().Errorw("delete lineup fail", "err", err)
		return err
	}

	return nil
}
