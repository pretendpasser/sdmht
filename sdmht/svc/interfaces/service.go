package itfs

import (
	"context"

	"sdmht/sdmht/svc/entity"
)

type Service interface{}
type SignalingService interface {
	NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRsp, error)

	KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error
	Offline(ctx context.Context, req *entity.LogoutReq) error
}
