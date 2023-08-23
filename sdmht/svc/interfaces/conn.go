package itfs

import "context"

type User2ConnRepo interface {
	Find(ctx context.Context, id []uint64) (map[uint64]string, error)
	Get(ctx context.Context, id uint64) (string, error)
	Add(ctx context.Context, id uint64, connName string) error
	Delete(ctx context.Context, id uint64) error
	FindConnNames(ctx context.Context) ([]string, error)
}
