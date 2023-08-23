package signaling_grpc

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht/api"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
)

var _ pb.SignalingServer = (*grpcServer)(nil)

// grpcServer is the Go kit server implementation for Signaling service.
type grpcServer struct {
	*pb.UnimplementedSignalingServer

	NewMatchHandler  grpc.Handler
	KeepAliveHandler grpc.Handler
	OfflineHandler   grpc.Handler
}

func NewGRPCServer(svc itfs.SignalingService, opts *kitx.ServerOptions) pb.SignalingServer {
	srv := &grpcServer{}

	logger := opts.Logger()
	tracer := opts.ZipkinTracer()

	options := []grpc.ServerOption{
		grpc.ServerBefore(opts.MetadataToCtx("sdmht")),
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	if tracer != nil {
		options = append(options, zipkin.GRPCServerTrace(tracer))
	}

	srv.NewMatchHandler = makeNewMatchHandler(svc, options, opts)
	srv.KeepAliveHandler = makeKeepAliveHandler(svc, options, opts)
	srv.OfflineHandler = makeOfflineHandler(svc, options, opts)

	return srv
}

func (s *grpcServer) NewMatch(ctx context.Context, req *pb.NewMatchReq) (*pb.NewMatchReply, error) {
	_, resp, err := s.NewMatchHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.NewMatchReply), nil
}

func (s *grpcServer) KeepAlive(ctx context.Context, req *pb.KeepAliveReq) (*pb.CommonReply, error) {
	_, resp, err := s.KeepAliveHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CommonReply), nil
}

func (s *grpcServer) Offline(ctx context.Context, req *pb.LogoutReq) (*pb.CommonReply, error) {
	_, resp, err := s.OfflineHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CommonReply), nil
}

func makeNewMatchHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeNewMatchEndpoint(svc)
		return ep, "webinar.signaling.NewMatch"
	}, opts)

	return grpc.NewServer(ep, decodeNewMatchReq, encodeNewMatchReply, options...)
}

func makeKeepAliveHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeKeepAliveEndpoint(svc)
		return ep, "webinar.signaling.KeepAlive"
	}, opts)

	return grpc.NewServer(ep, decodeKeepAliveReq, encodeCommonReply, options...)
}

func makeOfflineHandler(svc itfs.SignalingService, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeOfflineEndpoint(svc)
		return ep, "webinar.signaling.Offline"
	}, opts)

	return grpc.NewServer(ep, decodeLogoutReq, encodeCommonReply, options...)
}
