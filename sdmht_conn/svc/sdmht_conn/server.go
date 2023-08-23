package sdmht_conn

import (
	"context"
	"sync"
	"time"

	account_itfs "sdmht/account/svc/interfaces"
	"sdmht/lib"
	"sdmht/lib/log"
	sdmht_entity "sdmht/sdmht/svc/entity"
	"sdmht/sdmht_conn/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"
)

var _ itfs.ConnService = (*Server)(nil)

type Server struct {
	clients map[uint64]*Client // {catonID: *client}
	rwLock  sync.RWMutex

	ClientEventChan chan ClientEvent

	ConnMng    itfs.ConnManager
	AccountSvc account_itfs.Service

	closeChan chan struct{}

	ServerStartTime time.Time // 起服时间
}

func NewServer(connMng itfs.ConnManager, accountSvc account_itfs.Service) *Server {
	return &Server{
		clients:         make(map[uint64]*Client),
		ClientEventChan: make(chan ClientEvent, 10),
		ConnMng:         connMng,
		AccountSvc:      accountSvc,
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

func (s *Server) GetClient(catonID uint64) (*Client, bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	c, ok := s.clients[catonID]
	return c, ok
}

func (s *Server) AddClient(catonID uint64, client *Client) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if oldClient, ok := s.clients[catonID]; ok {
		if oldClient.conn != client.conn { // 防止客户端连上来后 反复发登录包
			s.tryKickClient(catonID)
		}
	}
	s.clients[catonID] = client
}

func (s *Server) RemoveClient(catonID uint64) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.removeClient(catonID)
}

func (s *Server) removeClient(catonID uint64) {
	delete(s.clients, catonID)
}

func (s *Server) tryKickClient(catonID uint64) {
	if client, ok := s.clients[catonID]; ok {
		client.CloseWait(entity.ClientQuitReasonKick)
		s.removeClient(catonID)
	}
}

func (s *Server) DispatchEventToClient(ctx context.Context, accountID uint64, event sdmht_entity.ClientEvent) (res sdmht_entity.DispatchEventToClientReply, err error) {
	log.S().Infow("handle dispatch", "accountID", accountID, "event.Type", event.Type, "event.Content", string(event.Content))
	c, has := s.GetClient(accountID)
	if !has {
		return res, lib.NewError(lib.ErrUnavailable, entity.ErrPocConnClientNotOnline)
	}
	payload := entity.NewReqPayload(c.NewSN(), event.Type, s.clients[accountID].Token(), event.Content)

	resp, err := c.DoRequest(payload)
	if resp.Result.Code != entity.ErrCodeMsgSuccess {
		res.OK = false
		res.ClientErr = resp.Result.Reason
	} else {
		res = sdmht_entity.DispatchEventToClientReply{
			UserID: c.AccountID(),
			OK:     true,
		}
	}
	return
}

// func (s *Server) SwitchToSpeecher(ctx context.Context, req webinar_entity.SwitchSpeecherReq) (res *webinar_entity.SwitchSpeecherRes, err error) {
// 	log.S().Infow("switch speeker", "operator", req.Operator, "new speaker", req.NewSpeecherID)
// 	c, has := s.GetClient(req.NewSpeecherID)
// 	if !has {
// 		log.S().Error("client not found")
// 		return res, lib.NewError(lib.ErrUnavailable, entity.ErrPocConnClientNotOnline)
// 	}
// 	payload := entity.NewReqPayload(c.NewSN(), webinar_entity.MsgTypeSwitchSpeecherCommand, "", req)
// 	resp, err := c.DoRequest(payload)
// 	if err != nil || resp.Code != entity.ErrCodeMsgSuccess {
// 		log.S().Errorw("client do request fail", "code", resp.Code, "reason", resp.Reason)
// 		return nil, lib.NewError(lib.ErrUnavailable, resp.Reason)
// 	}

// 	res = resp.MsgContent.(*webinar_entity.SwitchSpeecherRes)
// 	return
// }

func (s *Server) KickClient(ctx context.Context, catonID uint64) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.tryKickClient(catonID)
	return nil
}

func (s *Server) Close() {
	close(s.closeChan)
}
