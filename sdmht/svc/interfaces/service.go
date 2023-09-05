package itfs

import (
	"context"

	"sdmht/sdmht/svc/entity"
)

type Service interface {
	CreateLineup(ctx context.Context, lineup *entity.Lineup) error
	FindLineup(ctx context.Context, query *entity.LineupQuery) (int, []*entity.Lineup, error)
	GetLineup(ctx context.Context, id uint64, accountID uint64) (*entity.Lineup, error)
	UpdateLineup(ctx context.Context, lineup *entity.Lineup) error
	DeleteLineup(ctx context.Context, id uint64, accountID uint64) error
}
type SignalingService interface {
	Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error)

	KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error
	Offline(ctx context.Context, req *entity.LogoutReq) error

	NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error)
}
