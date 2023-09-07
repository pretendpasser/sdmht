package server

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht/api"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
)

var _ pb.SignalingServer = (*grpcServer)(nil)

// grpcServer is the Go kit server implementation for Signaling service.
type grpcServer struct {
	*pb.UnimplementedSignalingServer

	LoginHandler        grpc.Handler
	NewLineupHandler    grpc.Handler
	FindLineupHandler   grpc.Handler
	UpdateLineupHandler grpc.Handler
	DeleteLineupHandler grpc.Handler
	NewMatchHandler     grpc.Handler
	KeepAliveHandler    grpc.Handler
	OfflineHandler      grpc.Handler
}

func NewGRPCServer(svc itfs.SignalingService, opts *kitx.ServerOptions) pb.SignalingServer {
	srv := &grpcServer{}

	logger := opts.Logger()

	options := []grpc.ServerOption{
		grpc.ServerBefore(opts.MetadataToCtx("sdmht")),
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	srv.LoginHandler = makeLoginHandler(svc, options, opts)
	srv.NewLineupHandler = makeNewLineupHandler(svc, options, opts)
	srv.FindLineupHandler = makeFindLineupHandler(svc, options, opts)
	srv.UpdateLineupHandler = makeUpdateLineupHandler(svc, options, opts)
	srv.DeleteLineupHandler = makeDeleteLineupHandler(svc, options, opts)
	srv.NewMatchHandler = makeNewMatchHandler(svc, options, opts)
	srv.KeepAliveHandler = makeKeepAliveHandler(svc, options, opts)
	srv.OfflineHandler = makeOfflineHandler(svc, options, opts)

	return srv
}

func (s *grpcServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	_, res, err := s.LoginHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.LoginReply), nil
}

func (s *grpcServer) NewMatch(ctx context.Context, req *pb.NewMatchReq) (*pb.NewMatchReply, error) {
	_, res, err := s.NewMatchHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.NewMatchReply), nil
}

func (s *grpcServer) KeepAlive(ctx context.Context, req *pb.KeepAliveReq) (*pb.CommonReply, error) {
	_, res, err := s.KeepAliveHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CommonReply), nil
}

func (s *grpcServer) Offline(ctx context.Context, req *pb.LogoutReq) (*pb.CommonReply, error) {
	_, res, err := s.OfflineHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CommonReply), nil
}

func makeLoginHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeLoginEndpoint(svc)
		return ep, "sdmht.signaling.Login"
	}, opts)

	return grpc.NewServer(ep, deLoginReq, enLoginReply, options...)
}

func makeNewLineupHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeNewLineupEndpoint(svc)
		return ep, "sdmht.signaling.NewLineup"
	}, opts)

	return grpc.NewServer(ep, deNewLineupReq, enCommonReply, options...)
}

func makeFindLineupHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeFindLineupEndpoint(svc)
		return ep, "sdmht.signaling.FindLineup"
	}, opts)

	return grpc.NewServer(ep, deFindLineupReq, enFindLineupReply, options...)
}

func makeUpdateLineupHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeUpdateLineupEndpoint(svc)
		return ep, "sdmht.signaling.UpdateLineup"
	}, opts)

	return grpc.NewServer(ep, deUpdateLineupReq, enCommonReply, options...)
}

func makeDeleteLineupHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeDeleteLineupEndpoint(svc)
		return ep, "sdmht.signaling.DeleteLineup"
	}, opts)

	return grpc.NewServer(ep, deDeleteLineupReq, enCommonReply, options...)
}

func makeNewMatchHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeNewMatchEndpoint(svc)
		return ep, "sdmht.signaling.NewMatch"
	}, opts)

	return grpc.NewServer(ep, deNewMatchReq, enNewMatchReply, options...)
}

func makeKeepAliveHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeKeepAliveEndpoint(svc)
		return ep, "sdmht.signaling.KeepAlive"
	}, opts)

	return grpc.NewServer(ep, deKeepAliveReq, enCommonReply, options...)
}

func makeOfflineHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeOfflineEndpoint(svc)
		return ep, "sdmht.signaling.Offline"
	}, opts)

	return grpc.NewServer(ep, deLogoutReq, enCommonReply, options...)
}
