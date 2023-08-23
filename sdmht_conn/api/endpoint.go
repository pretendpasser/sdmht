package api

import (
	"context"

	"sdmht/lib/kitx"
	webinar_entity "sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
)

type Response = kitx.Response

func MakeDispatchEventToClientEndpoint(s itfs.ConnService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(webinar_entity.ClientEvent)
		rsp, err := s.DispatchEventToClient(ctx, req.UserID, req)
		return Response{Value: rsp, Error: err}, nil
	}
}

func MakeKickClientEndpoint(s itfs.ConnService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		termID := request.(uint64)
		_ = s.KickClient(ctx, termID)
		return Response{}, nil
	}
}
