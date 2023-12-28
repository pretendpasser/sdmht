package itfs

import (
	"context"
	"sdmht/sdmht_conn/svc/entity"
)

type SignalingService interface {
	Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error)
	NewLineup(ctx context.Context, req *entity.NewLineupReq) error
	FindLineup(ctx context.Context, req *entity.FindLineupReq) (*entity.FindLineupRes, error)
	UpdateLineup(ctx context.Context, req *entity.UpdateLineupReq) error
	DeleteLineup(ctx context.Context, req *entity.DeleteLineupReq) error

	KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error
	Offline(ctx context.Context, req *entity.LogoutReq) error

	NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error)
	JoinMatch(ctx context.Context, req *entity.JoinMatchReq) (*entity.JoinMatchRes, error)
	GetMatch(ctx context.Context, req *entity.GetMatchReq) (*entity.GetMatchRes, error)

	SyncOperate(ctx context.Context, req *entity.SyncOperateReq) (*entity.SyncOperateRes, error)
}
