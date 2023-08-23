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
	authenticateEndpoint endpoint.Endpoint
}

func (c *grpcClient) Authenticate(ctx context.Context, token string) (*entity.Account, error) {
	rsp, err := c.authenticateEndpoint(ctx, token)
	if err != nil {
		return nil, err
	}
	return rsp.(*entity.Account), nil
}

func NewClient(instancer sd.Instancer, opts *kitx.ClientOptions) itfs.Service {
	c := &grpcClient{}
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opts.SourceToGRPC()),
	}

	var serviceName = "account_pb.Account"

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

	return c
}
