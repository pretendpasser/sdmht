package api

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
)

func MakeNewMatchEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.NewMatchReq)
		rsp, err := s.NewMatch(ctx, req)
		return kitx.Response{Value: rsp, Error: err}, nil
	}
}

func MakeKeepAliveEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.KeepAliveReq)
		err = s.KeepAlive(ctx, req)
		return kitx.Response{Error: err}, nil
	}
}

func MakeOfflineEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.LogoutReq)
		err = s.Offline(ctx, req)
		return kitx.Response{Error: err}, nil
	}
}
