package server

import (
	"net"
	"sdmht/lib/log"
	"strings"
	"sync"
	"time"
)

type ConnFactory func(c net.Conn, s Server) Conn

var _ Server = (*server)(nil)

type heartBeat struct {
	timeout    time.Duration
	timeoutCnt int32
}

type server struct {
	listener   net.Listener
	conns      []Conn // key is neid
	heartBeats []*heartBeat
	connMu     sync.RWMutex
	// handler    ConnHandler
	closeC  chan struct{}
	factory ConnFactory
}

// func NewServer(listener net.Listener, factory ConnFactory, handler ConnHandler) *server {
func NewServer(listener net.Listener, factory ConnFactory) *server {
	return &server{
		listener: listener,
		// handler:    handler,
		closeC:     make(chan struct{}, 1),
		factory:    factory,
		conns:      make([]Conn, 128),
		heartBeats: make([]*heartBeat, 128),
	}
}

func (s *server) Serve() error {
	defer func() {
		log.S().Infow("server", "quit", "Serve")
	}()

	for {
		conn, err := s.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.S().Error("accept", "err", err)
			}
			return err
		}
		go func() {
			_ = conn.Serve()
			log.S().Warnw("neconn serve quit", "conn", conn)
		}()
	}
}

func (s *server) Accept() (Conn, error) {
	c, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	log.S().Infow("accept", "addr", c.RemoteAddr())

	conn := s.factory(c, s)
	return conn, nil
}

func (s *server) Close() error {
	s.closeC <- struct{}{}
	err := s.listener.Close() // close server listener first
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for _, conn := range s.conns { //  close all conns
		if connErr := conn.Close(); connErr != nil {
			log.S().Errorw("close conn", "err", connErr)
		}
	}
	return err
}
