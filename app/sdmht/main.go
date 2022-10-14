package main

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"sdmht/lib/log"
	"sdmht/lib/utils"
	"sdmht/svc/server"

	"go.uber.org/zap"
)

const (
	TimeoutDuration = 7 * time.Second
)

func main() {
	errc := make(chan error)
	var err error

	var (
		httpAddr       = utils.GetEnvDefault("HTTP_ADDR", ":10709")
		clientConnAddr = utils.GetEnvDefault("CLIENT_CONN_LISTEN_ADDR", ":11709")
	)

	listener, err := net.Listen("tcp", clientConnAddr)
	if err != nil {
		log.S().Panic(err)
	}

	var (
		connSrv = server.NewServer(listener, func(c net.Conn, s server.Server) server.Conn {
			return server.NewConn(c, s)
		})
	)

	httpSrv := &http.Server{Addr: httpAddr}

	go func() {
		errc <- httpSrv.ListenAndServe()
	}()

	go func() {
		errc <- connSrv.Serve()
	}()

	log.S().Info("run")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	select {
	case <-ctx.Done():
		log.S().Info("recv quit signal")
	case err = <-errc:
		log.S().Errorw("quit", "err", err)
	}
	closeContext, cancel := context.WithTimeout(context.TODO(), TimeoutDuration)
	defer cancel()
	if err = httpSrv.Shutdown(closeContext); err != nil {
		log.S().Errorw("httpSrv shutdown", zap.String("err", err.Error()))
	}
	if err = connSrv.Close(); err != nil {
		log.S().Errorw("connSrv shutdown", zap.String("err", err.Error()))
	}
	log.S().Info("exit success")
}
