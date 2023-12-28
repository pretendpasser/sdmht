package sdmht

import (
	"context"
	"sync"
	"time"

	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/sdmht_conn/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"
)

var _ itfs.ConnService = (*Server)(nil)

type Server struct {
	clients map[uint64]*Client // {accountID: *client}
	rwLock  sync.RWMutex

	ConnMng         itfs.SignalingService
	ServerStartTime time.Time // 起服时间
	ClientEventChan chan ClientEvent
	closeChan       chan struct{}
}

func NewServer(connMng itfs.SignalingService) *Server {
	return &Server{
		clients:         make(map[uint64]*Client),
		ClientEventChan: make(chan ClientEvent, 10),
		ConnMng:         connMng,
		closeChan:       make(chan struct{}, 1),
		ServerStartTime: time.Now(),
	}
}

func (s *Server) CountClients() int {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return len(s.clients)
}

func (s *Server) Serve() {
	onlineTicker := time.NewTicker(30 * time.Minute)
	defer onlineTicker.Stop()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.closeChan:
			return
		case event := <-s.ClientEventChan:
			if event.Type == ClientEventTypeAdd {
				log.S().Debugw("add client start", "clientID", event.AccountID, "event_addr", event.Client)
				s.AddClient(event.AccountID, event.Client)
				log.S().Debugw("add client end", "clientID", event.AccountID, "event_addr", event.Client)
			} else {
				log.S().Debugw("remove client start", "clientID", event.AccountID, "event_addr", event.Client)
				s.RemoveClient(event.AccountID)
				log.S().Debugw("remove client end", "clientID", event.AccountID, "event_addr", event.Client)
			}
		case <-onlineTicker.C:
			onlineNum := s.CountClients()
			log.S().Infow("CountClients", "online num", onlineNum)
		}
	}
}

func (s *Server) GetClient(accountID uint64) (*Client, bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	c, ok := s.clients[accountID]
	return c, ok
}

func (s *Server) AddClient(accountID uint64, client *Client) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if oldClient, ok := s.clients[accountID]; ok {
		if oldClient.conn != client.conn { // 防止客户端连上来后 反复发登录包
			s.tryKickClient(accountID)
		}
	}
	s.clients[accountID] = client
}

func (s *Server) RemoveClient(accountID uint64) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.removeClient(accountID)
}

func (s *Server) removeClient(accountID uint64) {
	delete(s.clients, accountID)
}

func (s *Server) tryKickClient(accountID uint64) {
	if client, ok := s.clients[accountID]; ok {
		client.CloseWait(entity.ClientQuitReasonKick)
		s.removeClient(accountID)
	}
}

func (s *Server) DispatchEventToClient(ctx context.Context, accountID uint64, event entity.ClientEvent) (res entity.DispatchEventToClientReply, err error) {
	log.S().Infow("handle dispatch", "accountID", accountID, "event.Type", event.Type, "event.Content", string(event.Content))
	c, has := s.GetClient(accountID)
	if !has {
		return res, lib.NewError(lib.ErrUnavailable, entity.ErrConnClientNotOnline)
	}
	payload := entity.NewReqPayload(c.NewSN(), event.Type, event.Content)

	resp, err := c.DoRequest(payload)
	if resp.Result.Code != lib.ErrSuccess {
		res.OK = false
		res.ClientErr = resp.Result.Reason
	} else {
		res = entity.DispatchEventToClientReply{
			AccountID: c.AccountID(),
			OK:        true,
		}
	}
	return
}

func (s *Server) KickClient(ctx context.Context, accountID uint64) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.tryKickClient(accountID)
	return nil
}

func (s *Server) Close() {
	close(s.closeChan)
}
