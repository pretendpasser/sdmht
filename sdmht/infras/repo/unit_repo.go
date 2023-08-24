package repo

import (
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.UnitRepo = (*unitRepo)(nil)

type unitRepo struct{}

func (s *unitRepo) Get() {

}
