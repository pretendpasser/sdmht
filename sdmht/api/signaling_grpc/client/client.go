package client

import (
	"context"

	"sdmht/lib/kitx"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ itfs.SignalingService = (*grpcClient)(nil)

// grpcClient is the Go kit client implementation for interfaces.Service.
type grpcClient struct {
	LoginEndpoint        endpoint.Endpoint
	NewLineupEndpoint    endpoint.Endpoint
	FindLineupEndpoint   endpoint.Endpoint
	UpdateLineupEndpoint endpoint.Endpoint
	DeleteLineupEndpoint endpoint.Endpoint
	NewMatchEndpoint     endpoint.Endpoint
	GetMatchEndpoint     endpoint.Endpoint
	JoinMatchEndpoint    endpoint.Endpoint
	SyncOperatorEndpoint endpoint.Endpoint
	KeepAliveEndpoint    endpoint.Endpoint
	OfflineEndpoint      endpoint.Endpoint
}

func (c *grpcClient) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	res, err := c.LoginEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.LoginRes), nil
}

func (c *grpcClient) NewLineup(ctx context.Context, req *entity.NewLineupReq) error {
	_, err := c.NewLineupEndpoint(ctx, req)
	return err
}

func (c *grpcClient) FindLineup(ctx context.Context, req *entity.FindLineupReq) (*entity.FindLineupRes, error) {
	res, err := c.FindLineupEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.FindLineupRes), nil
}

func (c *grpcClient) UpdateLineup(ctx context.Context, req *entity.UpdateLineupReq) error {
	_, err := c.UpdateLineupEndpoint(ctx, req)
	return err
}

func (c *grpcClient) DeleteLineup(ctx context.Context, req *entity.DeleteLineupReq) error {
	_, err := c.DeleteLineupEndpoint(ctx, req)
	return err
}

func (c *grpcClient) NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error) {
	res, err := c.NewMatchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.NewMatchRes), nil
}

func (c *grpcClient) GetMatch(ctx context.Context, req *entity.GetMatchReq) (*entity.GetMatchRes, error) {
	res, err := c.GetMatchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.GetMatchRes), nil
}

func (c *grpcClient) JoinMatch(ctx context.Context, req *entity.JoinMatchReq) (*entity.JoinMatchRes, error) {
	res, err := c.JoinMatchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.JoinMatchRes), nil
}

func (c *grpcClient) SyncOperator(ctx context.Context, req *entity.SyncOperator) error {
	_, err := c.SyncOperatorEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return nil
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

	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opts.MetadataToGRPC("sdmht")),
		grpctransport.ClientBefore(opts.SourceToGRPC()),
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
	c.NewLineupEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"NewLineup",
			enNewLineupReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.NewLineup"
	}, opts)
	c.FindLineupEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"FindLineup",
			enFindLineupReq,
			deFindLineupReply,
			pb.FindLineupReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.FindLineup"
	}, opts)
	c.UpdateLineupEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"UpdateLineup",
			enUpdateLineupReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.UpdateLineup"
	}, opts)
	c.DeleteLineupEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"DeleteLineup",
			enDeleteLineupReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.DeleteLineup"
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
	c.GetMatchEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"GetMatch",
			enGetMatchReq,
			deGetMatchReply,
			pb.GetMatchReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.GetMatch"
	}, opts)
	c.JoinMatchEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"JoinMatch",
			enJoinMatchReq,
			deJoinMatchReply,
			pb.JoinMatchReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.JoinMatch"
	}, opts)
	c.SyncOperatorEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"SyncOperator",
			enSyncOperatorReq,
			deCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "sdmht.signaling.rpc.SyncOperator"
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
