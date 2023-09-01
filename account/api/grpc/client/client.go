package client

import (
	"context"

	"sdmht/account/api/grpc/pb"
	"sdmht/account/svc/entity"
	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib/kitx"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ itfs.Service = (*grpcClient)(nil)

type grpcClient struct {
	registerEndpoint     endpoint.Endpoint
	loginEndpoint        endpoint.Endpoint
	logoutEndpoint       endpoint.Endpoint
	authenticateEndpoint endpoint.Endpoint
	getEndpoint          endpoint.Endpoint
}

func (c *grpcClient) Register(ctx context.Context, req *entity.RegisterReq) error {
	_, err := c.registerEndpoint(ctx, req)
	return err
}
func (c *grpcClient) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	res, err := c.loginEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*entity.LoginRes), nil
}
func (c *grpcClient) Logout(ctx context.Context, token string) error {
	_, err := c.logoutEndpoint(ctx, token)
	return err
}

func (c *grpcClient) Authenticate(ctx context.Context, token string) (*entity.Account, error) {
	res, err := c.authenticateEndpoint(ctx, token)
	if err != nil {
		return nil, err
	}
	return res.(*entity.Account), nil
}

func (c *grpcClient) GetAccount(ctx context.Context, id uint64) (*entity.Account, error) {
	res, err := c.getEndpoint(ctx, id)
	if err != nil {
		return nil, err
	}
	return res.(*entity.Account), nil
}

func NewClient(instancer sd.Instancer, opts *kitx.ClientOptions) itfs.Service {
	c := &grpcClient{}
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opts.SourceToGRPC()),
	}

	var serviceName = "account_pb.Account"

	c.registerEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Register",
			enRegisterReq,
			deRegisterRes,
			pb.RegisterRes{},
			options...,
		).Endpoint(), "account.rpc.Register"
	}, opts)

	c.loginEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Login",
			enLoginReq,
			deLoginRes,
			pb.LoginRes{},
			options...,
		).Endpoint(), "account.rpc.Login"
	}, opts)

	c.logoutEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Logout",
			enLogoutReq,
			deLogoutRes,
			pb.LogoutRes{},
			options...,
		).Endpoint(), "account.rpc.Logout"
	}, opts)

	c.authenticateEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Authenticate",
			enAuthenticateReq,
			deAuthenticateRes,
			pb.AuthenticateRes{},
			options...,
		).Endpoint(), "account.rpc.Authenticate"
	}, opts)

	c.getEndpoint = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"Get",
			enGetAccountReq,
			deGetAccountRes,
			pb.GetAccountRes{},
			options...,
		).Endpoint(), "account.rpc.Get"
	}, opts)

	return c
}
