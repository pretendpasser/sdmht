package account

import (
	"context"
	"sdmht/account/svc/entity"
	itfs "sdmht/account/svc/interfaces"
)

var _ itfs.Service = (*service)(nil)

type service struct {
}

func (s *service) Authenticate(context.Context, string) (*entity.Account, error) {
	return nil, nil
}
