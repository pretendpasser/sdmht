package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	accountcli "sdmht/account/api/grpc/client"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/seq"
	"sdmht/lib/utils"
	"sdmht/sdmht_conn/infras/repo"
	"sdmht/sdmht_conn/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"
	svc "sdmht/sdmht_conn/svc/sdmht"

	"github.com/go-kit/kit/sd"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

const (
// CloseTimeoutDuration = 7 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var server *svc.Server

func socketHandler(w http.ResponseWriter, r *http.Request) {
	heartBeatInterval, err := strconv.Atoi(utils.GetEnvDefault("HEARTBEAT_INTERVAL", "10"))
	if err != nil {
		log.S().Errorw("invalid heartBeatInterval", "error", err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.S().Errorw("something went wrong", "error", err)
		return
	}
	//defer conn.Close()
	log.S().Infow("new conn", "remote addr", conn.RemoteAddr().String())
	// TODO
	addr := utils.ClientIP(r)
	client := svc.NewClient(conn, addr, server.ClientEventChan, server.ConnMng, server.ServerStartTime, heartBeatInterval)
	go client.Run()
}

type config struct {
	LogLevel           string
	LogDisableSampling bool
	HttpAddr           string
	WSListenAddr       string
	AccountSvcAddr     string
	ServeAddr          string

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
	var err error
	errChan := make(chan error)

	cfg := new(config)
	parseConfig(cfg)
	data, _ := json.Marshal(cfg)
	log.S().Infow("cfg", "data", json.RawMessage(data))

	if err := initLog(cfg); err != nil {
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
		idGenerator   = seq.New()
		skillRepo     = entity.NewSkillList()
		lineupRepo    = repo.NewLineupRepo(db)
		unitRepo      = repo.NewUnitRepo(db, *skillRepo)
		matchRepo     = repo.NewMatchRepo("sdmht:match", redisDB)
		user2ConnRepo = repo.NewUser2ConnRepo("sdmht:user2conn", redisDB)
	)

	cliOpts := kitx.NewClientOptions(kitx.WithLogger(log.GetLogger()),
		kitx.WithLoadBalance(3, 5*time.Second), kitx.WithMetadata(map[string][]string{"sdmht_conn_addr": {cfg.ServeAddr}}))
	connMgr := svc.NewConnManager(user2ConnRepo, func(connName string) itfs.ConnService {
		return server
	})
	accountSvc := accountcli.NewClient(sd.FixedInstancer([]string{cfg.AccountSvcAddr}), cliOpts)
	signalingSvc := svc.NewSignalingService(idGenerator, lineupRepo, unitRepo, matchRepo, accountSvc, connMgr)

	server = svc.NewServer(signalingSvc)
	go server.Serve()
	http.HandleFunc("/sdmht", socketHandler)
	go func() {
		errChan <- http.ListenAndServe(cfg.WSListenAddr, nil)
	}()

	log.S().Info("run")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		log.S().Info("recv quit signal")
	case err = <-errChan:
		log.S().Errorw("quit", "err", err)
	}

	server.Close()

	log.S().Info("quit")
}

func parseConfig(config *config) {
	var err error

	config.LogLevel = utils.GetEnvDefault("LOG_LEVEL", "info")
	logDisableSamplingStr := utils.GetEnvDefault("LOG_DISABLE_SAMPLING", "false")
	if logDisableSamplingStr == "true" {
		config.LogDisableSampling = true
	}
	config.WSListenAddr = utils.GetEnvDefault("WS_LISTEN_ADDR", ":4090")
	config.HttpAddr = utils.GetEnvDefault("HTTP_LISTEN_ADDR", ":8090")
	config.AccountSvcAddr = utils.GetEnvDefault("ACCOUNT_ACCESS_ADDR", ":7001")
	config.ServeAddr = utils.GetEnvDefault("SERVE_ADDR", ":7090")
	if config.ServeAddr == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.S().Panic(err)
		}
		port := strings.Split(config.HttpAddr, ":")
		config.ServeAddr = fmt.Sprintf("%s:%s", hostname, port[1])
	}
	log.S().Infof("serveAddr: %s", config.ServeAddr)

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
	config.DBUrl = utils.GetEnvDefault("DB_URL", "root:123456@tcp(db:3306)/sdmht?parseTime=true&multiStatements=true")
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

func initLog(cfg *config) error {
	return log.InitLogger(log.WithLevel(cfg.LogLevel), log.DisableStacktrace(cfg.LogDisableSampling))
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
