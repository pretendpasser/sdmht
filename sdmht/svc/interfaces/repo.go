package itfs

import (
	"context"

	"sdmht/sdmht/svc/entity"
)

type UnitRepo interface {
	Get(ctx context.Context, ids []uint64) ([]*entity.Unit, error)
	Find(ctx context.Context, query *entity.UnitQuery) (int, []*entity.Unit, error)
}

type MatchRepo interface {
	New(match *entity.Match) error
	Set(match *entity.Match)
	Get(id uint64) (*entity.Match, error)
	Delete(id uint64) error

	RAdd(ctx context.Context, match *entity.Match) error
	RGet(ctx context.Context, id uint64) (*entity.Match, error)
	RUpdate(ctx context.Context, match *entity.Match) error
	RDelete(ctx context.Context, id uint64) error
}

type LineupRepo interface {
	Get(ctx context.Context, accountID uint64, id uint64) (*entity.Lineup, error)
	Find(ctx context.Context, query *entity.LineupQuery) (total int, res []*entity.Lineup, err error)
	Create(ctx context.Context, lineup *entity.Lineup) error
	Update(ctx context.Context, lineup *entity.Lineup) error
	Delete(ctx context.Context, accountID uint64, id uint64) error
}
