package client

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/lib/log"
	sdmht_entity "sdmht/sdmht/svc/entity"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"
	itfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ itfs.ConnService = (*grpcClient)(nil)

type grpcClient struct {
	dispatchEventToClient endpoint.Endpoint
	kickClient            endpoint.Endpoint
}

func NewClient(instancer sd.Instancer, opts *kitx.ClientOptions) *grpcClient {
	c := &grpcClient{}

	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opts.MetadataToGRPC("sdmht")),
		grpctransport.ClientBefore(jwt.ContextToGRPC()),
	}

	var serviceName = "conn_pb.Conn"
	c.dispatchEventToClient = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"DispatchEventToClient",
			encodeDispatchEventToClientRequest,
			decodeDispatchEventToClientReply,
			pb.DispatchEventToClientReply{},
			options...,
		).Endpoint(), "Conn.rpc.DispatchEventToClient"
	}, opts)

	c.kickClient = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			serviceName,
			"KickClient",
			encodeKickClientRequest,
			decodeCommonReply,
			pb.CommonReply{},
			options...,
		).Endpoint(), "Conn.rpc.KickClient"
	}, opts)

	return c
}

func (c *grpcClient) DispatchEventToClient(ctx context.Context, accountID uint64,
	event sdmht_entity.ClientEvent) (sdmht_entity.DispatchEventToClientReply, error) {
	log.S().Infow("DispatchEventToClient", "accountID", accountID, "req", event)
	rsp, err := c.dispatchEventToClient(ctx, &event)
	if err != nil {
		return sdmht_entity.DispatchEventToClientReply{}, err
	}
	return rsp.(sdmht_entity.DispatchEventToClientReply), nil
}

func (c *grpcClient) KickClient(ctx context.Context, termID uint64) error {
	_, err := c.kickClient(ctx, termID)
	return err
}
