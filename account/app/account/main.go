package main

import (
	"context"
	"encoding/json"
	"net"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	account_pb "sdmht/account/api/grpc/pb"
	account_grpc "sdmht/account/api/grpc/server"
	"sdmht/account/infras/repo"
	svc "sdmht/account/svc/account"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/utils"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type config struct {
	LogLevel           string
	LogDisableSampling bool
	HttpAddr           string
	GrpcAddr           string
	// redis
	RedisDBUrl        string
	RedisMaxIdleTime  int // second
	RedisMaxLifeTime  int
	RedisMaxIdleConns int
	RedisMaxOpenConns int
	// mysql
	DBUrl               string
	MysqlDBMaxIdleTime  int
	MysqlDBMaxLifeTime  int
	MysqlDBMaxIdleConns int
	MysqlDBMaxOpenConns int
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

	db, err := initMysqlDB(cfg.DBUrl, cfg.MysqlDBMaxIdleTime, cfg.MysqlDBMaxLifeTime,
		cfg.MysqlDBMaxIdleConns, cfg.MysqlDBMaxOpenConns)
	if err != nil {
		log.S().Panic("err", err)
	}
	defer db.Close()

	redisDB, err := initRedis(cfg.RedisDBUrl, cfg.RedisMaxIdleTime, cfg.RedisMaxLifeTime,
		cfg.RedisMaxIdleConns, cfg.RedisMaxOpenConns)
	if err != nil {
		log.S().Panic("err", err)
	}
	defer redisDB.Close()

	var (
		accountRepo = repo.NewAccountRepo(db)
		tokenRepo   = repo.NewTokenRepo(redisDB, "sdmht:token:")
	)
	srvOpts := kitx.NewServerOptions(kitx.WithLogger(log.GetLogger()), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(nil))
	// cliOpts := kitx.NewClientOptions(kitx.WithLogger(log.GetLogger()), kitx.WithLoadBalance(3, 5*time.Second), kitx.WithZipkinTracer(nil))

	accountSvc := svc.NewService(accountRepo, tokenRepo)
	grpcServer := account_grpc.NewGRPCServer(accountSvc, srvOpts)
	grpcService := grpc.NewServer()
	account_pb.RegisterAccountServer(grpcService, grpcServer)
	go func() {
		grpcListener, err := net.Listen("tcp", cfg.GrpcAddr)
		if err != nil {
			errc <- err
			return
		}
		errc <- grpcService.Serve(grpcListener)
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

	// mysql
	config.DBUrl = utils.GetEnvDefault("DB_URL", "root:123456@tcp(c3_db:3306)/c3_sla?parseTime=true&multiStatements=true")
	if dbMaxIdleTimeStr := utils.GetEnvDefault("MYSQL_DB_MAX_IDLE_TIME", ""); dbMaxIdleTimeStr != "" {
		config.MysqlDBMaxIdleTime, err = strconv.Atoi(dbMaxIdleTimeStr)
		if err != nil {
			log.S().Panic("MYSQLDB_MAX_IDLE_TIME invalid", dbMaxIdleTimeStr)
		}
	}
	if dbMaxLifeTimeStr := utils.GetEnvDefault("MYSQL_DB_MAX_LIFE_TIME", ""); dbMaxLifeTimeStr != "" {
		config.MysqlDBMaxIdleConns, err = strconv.Atoi(dbMaxLifeTimeStr)
		if err != nil {
			log.S().Panic("MYSQLDB_MAX_IDLE_CONNS invalid", dbMaxLifeTimeStr)
		}
	}
	if dbMaxIdleConnsStr := utils.GetEnvDefault("MYSQL_DB_MAX_IDLE_CONNS", ""); dbMaxIdleConnsStr != "" {
		config.MysqlDBMaxIdleConns, err = strconv.Atoi(dbMaxIdleConnsStr)
		if err != nil {
			log.S().Panic("MYSQLDB_MAX_IDLE_CONNS invalid", dbMaxIdleConnsStr)
		}
	}
	if dbMaxOpenConnsStr := utils.GetEnvDefault("MYSQL_DB_MAX_OPEN_CONNS", ""); dbMaxOpenConnsStr != "" {
		config.MysqlDBMaxOpenConns, err = strconv.Atoi(dbMaxOpenConnsStr)
		if err != nil {
			log.S().Panic("MYSQLDB_MAX_OPEN_CONNS invalid", dbMaxOpenConnsStr)
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

func initMysqlDB(dbUrl string, dbMaxIdleTime int, dbMaxLifeTime int, dbMaxIdleConns int, dbMaxOpenConns int) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}

	if dbMaxIdleTime != 0 {
		db.DB.SetConnMaxIdleTime(time.Duration(dbMaxIdleTime) * time.Second)
	}
	if dbMaxLifeTime != 0 {
		db.DB.SetConnMaxLifetime(time.Duration(dbMaxLifeTime) * time.Second)
	}
	if dbMaxIdleConns != 0 {
		db.DB.SetMaxIdleConns(dbMaxIdleConns)
	}

	if dbMaxOpenConns != 0 {
		db.DB.SetMaxOpenConns(dbMaxOpenConns)
	}

	return db, nil
}
