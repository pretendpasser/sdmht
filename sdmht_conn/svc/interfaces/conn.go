package itfs

import (
	"context"
	"sdmht/sdmht_conn/svc/entity"
)

type User2ConnRepo interface {
	Find(ctx context.Context, id []uint64) (map[uint64]string, error)
	Get(ctx context.Context, id uint64) (string, error)
	Add(ctx context.Context, id uint64, wechatID string) error
	Delete(ctx context.Context, id uint64) error
	FindConnNames(ctx context.Context) ([]string, error)
}

type ConnService interface {
	DispatchEventToClient(ctx context.Context, target uint64, event entity.ClientEvent) (entity.DispatchEventToClientReply, error)
	KickClient(ctx context.Context, id uint64) error
}
