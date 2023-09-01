package api

import (
	"context"

	"sdmht/account/svc/entity"
	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib/kitx"

	"github.com/go-kit/kit/endpoint"
)

func MakeRegisterEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.RegisterReq)
		err = s.Register(ctx, req)
		return kitx.Response{Value: nil, Error: err}, nil
	}
}

func MakeLoginEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.LoginReq)
		res, err := s.Login(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeLogoutEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := request.(string)
		err = s.Logout(ctx, token)
		return kitx.Response{Value: nil, Error: err}, nil
	}
}

func MakeAuthenticateEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := request.(string)
		v, err := s.Authenticate(ctx, token)
		return kitx.Response{Value: v, Error: err}, nil
	}
}

func MakeGetAccountEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(uint64)
		res, err := s.GetAccount(ctx, id)
		return kitx.Response{Value: res, Error: err}, nil
	}
}
