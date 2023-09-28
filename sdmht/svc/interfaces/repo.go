package itfs

import (
	"context"

	"sdmht/sdmht/svc/entity"
)

type UnitRepo interface {
	Get(ctx context.Context, ids []int64) ([]*entity.Unit, error)
	Find(ctx context.Context, query *entity.UnitQuery) (int, []*entity.Unit, error)
}

type MatchRepo interface {
	SetByAccount(accountID uint64, matchID uint64)
	GetByAccount(accountID uint64) (matchID uint64)

	New(match *entity.Match) error
	Join(match *entity.Match) error
	Set(match *entity.Match)
	Get(id uint64) (entity.Match, error)
	Delete(id uint64)

	RSet(ctx context.Context, match *entity.Match) error
	RGet(ctx context.Context, accountID uint64) (uint64, error)
	RDelete(ctx context.Context, accountID uint64) error
	RHSet(ctx context.Context, match *entity.Match) error
	RHGet(ctx context.Context, id uint64) (*entity.Match, error)
	RHDelete(ctx context.Context, id uint64) error
}

type LineupRepo interface {
	Get(ctx context.Context, accountID uint64, id uint64) (*entity.Lineup, error)
	Find(ctx context.Context, query *entity.LineupQuery) (total int, res []*entity.Lineup, err error)
	Create(ctx context.Context, lineup *entity.Lineup) error
	Update(ctx context.Context, lineup *entity.Lineup) error
	Delete(ctx context.Context, accountID uint64, id uint64) error
}
