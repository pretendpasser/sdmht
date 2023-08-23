package grpc

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht_conn/api"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"
	itfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

var _ pb.ConnServer = (*grpcServer)(nil)

type grpcServer struct {
	dispatchEventToClient grpctransport.Handler
	kickClient            grpctransport.Handler

	pb.UnimplementedConnServer
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

func NewGRPCServer(svc itfs.ConnService, opts *kitx.ServerOptions) pb.ConnServer {
	srv := &grpcServer{}

	logger := opts.Logger()
	tracer := opts.ZipkinTracer()

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if tracer != nil {
		options = append(options, zipkin.GRPCServerTrace(tracer))
	}

	srv.dispatchEventToClient = makeDispatchEventToClientHandler(svc, options, opts)
	srv.kickClient = makeKickClientHandler(svc, options, opts)

	return srv
}

func makeDispatchEventToClientHandler(svc itfs.ConnService, options []grpctransport.ServerOption,
	opts *kitx.ServerOptions) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeDispatchEventToClientEndpoint(svc)
		return ep, "sdmht_conn.DispatchEventToClient"
	}, opts)

	return grpctransport.NewServer(ep, decodeDispatchEventToClientReq, encodeDispatchEventToClientReply, options...)
}

func makeKickClientHandler(svc itfs.ConnService, options []grpctransport.ServerOption,
	opts *kitx.ServerOptions) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeKickClientEndpoint(svc)
		return ep, "sdmht_conn.KickClient"
	}, opts)

	return grpctransport.NewServer(ep, decodeKickClientReq, encodeKickClientReply, options...)
}
