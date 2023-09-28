package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/utils"
	sdmht_client "sdmht/sdmht/api/signaling_grpc/client"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"
	apigrpc "sdmht/sdmht_conn/api/grpc/server"
	"sdmht/sdmht_conn/svc/sdmht_conn"

	"github.com/go-kit/kit/sd"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

const (
// CloseTimeoutDuration = 7 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var server *sdmht_conn.Server

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
	client := sdmht_conn.NewClient(conn, addr, server.ClientEventChan, server.ConnMng, server.ServerStartTime, heartBeatInterval)
	go client.Run()
}

type config struct {
	WSListenAddr   string
	GrpcAddr       string
	ZipkinUrl      string
	ID             int // 与mng握手用
	MngAddr        string
	AccountSvcAddr string
	ServeAddr      string

	LogLevel           string
	LogDisableSampling bool
}

func parseConfig(config *config) {
	idStr := utils.GetEnvDefault("ID", "1")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.S().Panicw("id", "err", err)
	}
	config.ID = id
	config.WSListenAddr = utils.GetEnvDefault("WS_LISTEN_ADDR", ":4090")
	config.GrpcAddr = utils.GetEnvDefault("GRPC_LISTEN_ADDR", ":7091")
	config.MngAddr = utils.GetEnvDefault("MNG_LISTEN_ADDR", ":7090")
	config.AccountSvcAddr = utils.GetEnvDefault("ACCOUNT_ACCESS_ADDR", ":7001")
	config.ServeAddr = utils.GetEnvDefault("SERVE_ADDR", "")
	config.LogLevel = utils.GetEnvDefault("LOG_LEVEL", "info")
	logDisableSamplingStr := utils.GetEnvDefault("LOG_DISABLE_SAMPLING", "false")
	if logDisableSamplingStr == "true" {
		config.LogDisableSampling = true
	}
}

func main() {
	if err := initLog(); err != nil {
		panic(err)
	}

	cfg := new(config)
	parseConfig(cfg)

	data, _ := json.Marshal(cfg)
	log.S().Infow("cfg", "data", json.RawMessage(data))

	var err error
	errChan := make(chan error)

	if cfg.ServeAddr == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.S().Panic(err)
		}
		port := strings.Split(cfg.GrpcAddr, ":")
		cfg.ServeAddr = fmt.Sprintf("%s:%s", hostname, port[1])
	}
	log.S().Infof("serveAddr: %s", cfg.ServeAddr)

	cliOpts := kitx.NewClientOptions(kitx.WithLogger(log.GetLogger()),
		kitx.WithLoadBalance(3, 5*time.Second), kitx.WithMetadata(map[string][]string{"sdmht_conn_addr": {cfg.ServeAddr}}))

	sdmhtMng := sdmht_client.NewClient(sd.FixedInstancer([]string{cfg.MngAddr}), cliOpts)

	server = sdmht_conn.NewServer(sdmhtMng)
	srvOpts := kitx.NewServerOptions(
		kitx.WithLogger(log.GetLogger()),
		kitx.WithRateLimit(nil),
		kitx.WithCircuitBreaker(0),
		kitx.WithMetrics(nil),
	)

	grpcServer := grpc.NewServer()
	grpcService := apigrpc.NewGRPCServer(server, srvOpts)
	pb.RegisterConnServer(grpcServer, grpcService)
	grpcListener, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.S().Panic("grpc listen err", err)
	}
	go func() {
		errChan <- grpcServer.Serve(grpcListener)
	}()

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

	grpcServer.GracefulStop()

	log.S().Info("quit")
}

func initLog() error {
	level := utils.GetEnvDefault("LOG_LEVEL", "info")
	return log.InitLogger(log.WithLevel(level), log.DisableStacktrace(true))
}
