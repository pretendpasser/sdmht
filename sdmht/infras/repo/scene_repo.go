package repo

import (
	itfs "sdmht/sdmht/svc/interfaces"
)

var _ itfs.Scene = (*sceneRepo)(nil)

type sceneRepo struct{}

func (s *sceneRepo) Get() {

}
