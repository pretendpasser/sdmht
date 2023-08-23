package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	accountcli "sdmht/account/api/grpc/client"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/seq"
	"sdmht/lib/utils"
	apihttp "sdmht/sdmht/api/http"
	"sdmht/sdmht/api/signaling_grpc"
	signaling_pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/infras/repo"
	sdmht_svc "sdmht/sdmht/svc/sdmht"
	conncli "sdmht/sdmht_conn/api/grpc/client"
	connitfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/go-kit/kit/sd"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type config struct {
	LogLevel           string
	LogDisableSampling bool
	HttpAddr           string
	GrpcAddr           string

	// redis
	RedisMaxIdleTime  int // second
	RedisMaxLifeTime  int
	RedisMaxIdleConns int
	RedisMaxOpenConns int
	RedisDBUrl        string
}

func main() {
	errc := make(chan error)
	var err error

	cfg := new(config)
	parsesConfig(cfg)
	data, _ := json.Marshal(cfg)
	log.S().Infow("cfg", "data", json.RawMessage(data))

	if err := log.InitLogger(log.WithLevel(cfg.LogLevel), log.DisableStacktrace(cfg.LogDisableSampling)); err != nil {
		panic(err)
	}

	redisDB, err := initRedis(cfg.RedisDBUrl, cfg.RedisMaxIdleTime, cfg.RedisMaxLifeTime,
		cfg.RedisMaxIdleConns, cfg.RedisMaxOpenConns)
	if err != nil {
		log.S().Panic("err", err)
	}
	defer redisDB.Close()

	var (
		idGenerator = seq.New()
	)

	srvOpts := kitx.NewServerOptions(kitx.WithLogger(log.GetLogger()), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(nil))
	cliOpts := kitx.NewClientOptions(kitx.WithLogger(log.GetLogger()), kitx.WithLoadBalance(3, 5*time.Second), kitx.WithZipkinTracer(nil))

	user2ConnRepo := repo.NewUser2ConnRepo("sdmht:user2conn", redisDB)
	connMgr := sdmht_svc.NewConnManager(user2ConnRepo, func(connName string) connitfs.ConnService {
		instance := []string{connName}
		return conncli.NewClient(sd.FixedInstancer(instance), cliOpts)
	})
	sdmhtSvc := sdmht_svc.NewService()
	accountSvc := accountcli.NewClient(sd.FixedInstancer([]string{utils.GetEnvDefault("ACCOUNT_ACCESS_ADDR", "account:7001")}), cliOpts)
	signalingSvc := sdmht_svc.NewSignalingService(idGenerator, nil, accountSvc, connMgr)
	grpcServer := signaling_grpc.NewGRPCServer(signalingSvc, srvOpts)
	grpcService := grpc.NewServer()
	signaling_pb.RegisterSignalingServer(grpcService, grpcServer)
	go func() {
		grpcListener, err := net.Listen("tcp", cfg.GrpcAddr)
		if err != nil {
			errc <- err
			return
		}
		errc <- grpcService.Serve(grpcListener)
	}()

	go func() {
		errc <- http.ListenAndServe(cfg.HttpAddr, apihttp.NewHTTPHandler(sdmhtSvc, srvOpts))
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

func parsesConfig(config *config) {
	var err error

	config.LogLevel = utils.GetEnvDefault("LOG_LEVEL", "info")
	logDisableSamplingStr := utils.GetEnvDefault("LOG_DISABLE_SAMPLING", "false")
	if logDisableSamplingStr == "true" {
		config.LogDisableSampling = true
	}
	config.HttpAddr = utils.GetEnvDefault("HTTP_LISTEN_ADDR", ":8090")
	config.GrpcAddr = utils.GetEnvDefault("GRPC_LISTEN_ADDR", ":9090")

	// redis
	config.RedisDBUrl = utils.GetEnvDefault("REDIS_DB_URL", "redis://redis:6379")
	redisMaxIdleTimeStr := utils.GetEnvDefault("REDIS_MAX_IDLE_TIME", "")
	if redisMaxIdleTimeStr != "" {
		config.RedisMaxIdleTime, err = strconv.Atoi(redisMaxIdleTimeStr)
		if err != nil {
			log.S().Panic("REDIS_MAX_IDLE_TIME invalid", redisMaxIdleTimeStr)
		}
	}
	redisMaxLifeTimeStr := utils.GetEnvDefault("REDIS_MAX_LIFE_TIME", "")
	if redisMaxLifeTimeStr != "" {
		config.RedisMaxLifeTime, err = strconv.Atoi(redisMaxLifeTimeStr)
		if err != nil {
			log.S().Panic("REDIS_MAX_LIFE_TIME invalid", redisMaxLifeTimeStr)
		}
	}
	redisMaxIdleConnsStr := utils.GetEnvDefault("REDIS_MAX_IDLE_CONNS", "")
	if redisMaxIdleConnsStr != "" {
		config.RedisMaxIdleConns, err = strconv.Atoi(redisMaxIdleConnsStr)
		if err != nil {
			log.S().Panic("REDIS_MAX_IDLE_CONNS invalid", redisMaxIdleConnsStr)
		}
	}
	redisMaxOpenConnsStr := utils.GetEnvDefault("REDIS_MAX_OPEN_CONNS", "")
	if redisMaxOpenConnsStr != "" {
		config.RedisMaxOpenConns, err = strconv.Atoi(redisMaxOpenConnsStr)
		if err != nil {
			log.S().Panic("REDIS_MAX_OPEN_CONNS invalid", redisMaxOpenConnsStr)
		}
	}

}

func initRedis(redisUrl string, idleTimeout int, maxConnAge int, minIdelConns int, maxPoolSize int) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	if idleTimeout != 0 {
		opts.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
	if maxConnAge != 0 {
		opts.MaxConnAge = time.Duration(minIdelConns) * time.Second
	}
	if minIdelConns != 0 {
		opts.MinIdleConns = minIdelConns
	}
	if maxPoolSize != 0 {
		opts.PoolSize = maxPoolSize
	}

	return redis.NewClient(opts), nil
}