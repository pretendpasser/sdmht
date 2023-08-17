package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/metriccollector"
	"sdmht/lib/utils"
	apihttp "sdmht/sdmht/api/http"
	sdmht_svc "sdmht/sdmht/svc/sdmht"
)

func main() {
	errc := make(chan error)
	var err error
	log_level := utils.GetEnvDefault("LOG_LEVEL", "info")
	if err := log.InitLogger(log.WithLevel(log_level), log.DisableStacktrace(true)); err != nil {
		panic(err)
	}

	var (
		httpAddr = utils.GetEnvDefault("HTTP_ADDR", ":7709")
		// grpcAddr = utils.GetEnvDefault("GRPC_LISTEN_ADDR", ":8709")
	)

	metrics := metriccollector.NewRequestLatencyMetrics("sdmht")
	srvOpts := kitx.NewServerOptions(kitx.WithLogger(log.GetLogger()), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(metrics), kitx.WithZipkinTracer(nil))

	sdmhtSvc := sdmht_svc.NewService()

	// grpcServer := signaling_grpc.NewGRPCServer(signalingSvc, srvOpts)
	// grpcService := grpc.NewServer()
	// signaling_pb.RegisterSignalingServer(grpcService, grpcServer)
	// go func() {
	// 	grpcListener, err := net.Listen("tcp", grpcAddr)
	// 	if err != nil {
	// 		errc <- err
	// 		return
	// 	}
	// 	errc <- grpcService.Serve(grpcListener)
	// }()

	// httpSrv := &http.Server{Addr: httpAddr, Handler: apihttp.NewHTTPHandler(webinarSvc, srvOpts)}
	// go func() {
	// 	errc <- httpSrv.ListenAndServe()
	// }()

	go func() {
		errc <- http.ListenAndServe(httpAddr, apihttp.NewHTTPHandler(sdmhtSvc, srvOpts))
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
	log.S().Info("exit success")
}
