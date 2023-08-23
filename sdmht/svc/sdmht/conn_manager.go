package sdmht

import (
	"context"
	"sync"

	itfs "sdmht/sdmht/svc/interfaces"
	connitfs "sdmht/sdmht_conn/svc/interfaces"
)

type ConnSvcFactory func(connName string) connitfs.ConnService

type userConnSvcCache struct {
	cache map[string]connitfs.ConnService
	mu    sync.RWMutex
}

func newUserConnSvcCache() *userConnSvcCache {
	return &userConnSvcCache{
		cache: make(map[string]connitfs.ConnService),
	}
}

func (c *userConnSvcCache) load(serveAddr string, factory ConnSvcFactory) connitfs.ConnService {
	c.mu.RLock()
	svc, ok := c.cache[serveAddr]
	if ok {
		c.mu.RUnlock()
		return svc
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	svc, ok = c.cache[serveAddr] // Double check because other goroutines may create it first
	if !ok {
		svc = factory(serveAddr)
		c.cache[serveAddr] = svc
	}

	return svc
}

type ConnManager struct {
	uid2ConnRepo itfs.User2ConnRepo
	factory      ConnSvcFactory
	cache        *userConnSvcCache
}

func NewConnManager(uid2ConnRepo itfs.User2ConnRepo, factory ConnSvcFactory) *ConnManager {
	return &ConnManager{
		uid2ConnRepo: uid2ConnRepo,
		factory:      factory,
		cache:        newUserConnSvcCache(),
	}
}

func (m *ConnManager) GetConnCli(ctx context.Context, id uint64) (connitfs.ConnService, error) {
	addr, err := m.uid2ConnRepo.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	cli := m.cache.load(addr, m.factory)
	return cli, nil
}

func (m *ConnManager) User2ConnRepo() itfs.User2ConnRepo {
	return m.uid2ConnRepo
}
