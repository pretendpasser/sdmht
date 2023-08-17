package sdmht

import (
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.Service = (*service)(nil)

type service struct {
}

func NewService() *service {
	return &service{}
}
