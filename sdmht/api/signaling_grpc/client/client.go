package client

import (
	"context"

	"sdmht/lib/kitx"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ itfs.SignalingService = (*grpcClient)(nil)

// grpcClient is the Go kit client implementation for interfaces.Service.
type grpcClient struct {
	LoginEndpoint     endpoint.Endpoint
	NewMatchEndpoint  endpoint.Endpoint
	KeepAliveEndpoint endpoint.Endpoint
	OfflineEndpoint   endpoint.Endpoint
}

func (c *grpcClient) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	res, err := c.LoginEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.LoginRes), nil
}

func (c *grpcClient) NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error) {
	res, err := c.NewMatchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.NewMatchRes), nil
}

func (c *grpcClient) KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error {
	_, err := c.KeepAliveEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *grpcClient) Offline(ctx context.Context, req *entity.LogoutReq) error {
	_, err := c.OfflineEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func NewClient(instancer sd.Instancer, opts *kitx.ClientOptions) itfs.SignalingService {
	c := &grpcClient{}

	tracer := opts.ZipkinTracer()

	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opts.MetadataToGRPC("sdmht")),
		grpctransport.ClientBefore(opts.SourceToGRPC()),
	}
	if tracer != nil {
		options = append(options, zipkin.GRPCClientTrace(tracer))
	}

	var serviceName = "signaling_pb.Signaling"
	c.LoginEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Login",
			enLoginReq,
			deLoginReply,
			pb.LoginReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.Login"
	}, opts)
	c.NewMatchEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"NewMatch",
			enNewMatchReq,
			deNewMatchReply,
			pb.NewMatchReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.NewMatch"
	}, opts)
	c.KeepAliveEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"KeepAlive",
			enKeepAliveReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.KeepAlive"
	}, opts)
	c.OfflineEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Offline",
			enLogoutReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.Logout"
	}, opts)

	return c
}
