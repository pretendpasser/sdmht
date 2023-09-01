package itfs

import (
	"context"
	"sdmht/account/svc/entity"
)

type Service interface {
	Register(ctx context.Context, req *entity.RegisterReq) error
	Login(ctx context.Context, req *entity.LoginReq) (res *entity.LoginRes, err error)
	Logout(ctx context.Context, token string) error
	Authenticate(ctx context.Context, token string) (*entity.Account, error)

	GetAccount(ctx context.Context, id uint64) (*entity.Account, error)
}
