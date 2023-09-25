package server

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht_conn/api"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"
	itfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
)

var _ pb.ConnServer = (*grpcServer)(nil)

type grpcServer struct {
	*pb.UnimplementedConnServer

	dispatchEventToClient grpc.Handler
	kickClient            grpc.Handler
}

func NewGRPCServer(svc itfs.ConnService, opts *kitx.ServerOptions) pb.ConnServer {
	srv := &grpcServer{}

	logger := opts.Logger()

	options := []grpc.ServerOption{
		grpc.ServerBefore(opts.MetadataToCtx("sdmht")),
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	srv.dispatchEventToClient = makeDispatchEventToClientHandler(svc, options, opts)
	srv.kickClient = makeKickClientHandler(svc, options, opts)

	return srv
}

func (s *grpcServer) DispatchEventToClient(ctx context.Context, req *pb.ClientEventReq) (*pb.DispatchEventToClientReply,
	error) {
	_, rsp, err := s.dispatchEventToClient.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.(*pb.DispatchEventToClientReply), nil
}

func (s *grpcServer) KickClient(ctx context.Context, req *pb.KickClientReq) (*pb.CommonReply, error) {
	_, rsp, err := s.kickClient.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.(*pb.CommonReply), nil
}

func makeDispatchEventToClientHandler(svc itfs.ConnService, options []grpc.ServerOption,
	opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeDispatchEventToClientEndpoint(svc)
		return ep, "sdmht_conn.DispatchEventToClient"
	}, opts)

	return grpc.NewServer(ep, decodeDispatchEventToClientReq, encodeDispatchEventToClientReply, options...)
}

func makeKickClientHandler(svc itfs.ConnService, options []grpc.ServerOption,
	opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeKickClientEndpoint(svc)
		return ep, "sdmht_conn.KickClient"
	}, opts)

	return grpc.NewServer(ep, decodeKickClientReq, encodeKickClientReply, options...)
}
