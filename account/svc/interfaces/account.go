package itfs

import (
	"context"
	"sdmht/account/svc/entity"
)

type Service interface {
	Authenticate(context.Context, string) (*entity.Account, error)
}
