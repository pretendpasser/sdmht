package itfs

import (
	"context"
	"time"

	"sdmht/account/svc/entity"
)

type AccountRepo interface {
	Add(ctx context.Context, account *entity.Account) error
	Get(ctx context.Context, id uint64) (*entity.Account, error)
	GetByWechatID(ctx context.Context, wechatID string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uint64) error
}

type TokenRepo interface {
	Add(ctx context.Context, token string, accountID uint64, ttl time.Duration) error
	Get(ctx context.Context, token string, ttl time.Duration) (uint64, error)
	Delete(ctx context.Context, token string) error
}
