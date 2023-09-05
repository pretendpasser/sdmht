package api

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-kit/kit/endpoint"
)

func MakeLoginEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.LoginReq)
		res, err := s.Login(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeNewMatchEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.NewMatchReq)
		res, err := s.NewMatch(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
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
