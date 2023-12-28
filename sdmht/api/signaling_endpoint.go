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

func MakeNewLineupEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.NewLineupReq)
		err = s.NewLineup(ctx, req)
		return kitx.Response{Value: nil, Error: err}, nil
	}
}

func MakeFindLineupEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.FindLineupReq)
		res, err := s.FindLineup(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeUpdateLineupEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.UpdateLineupReq)
		err = s.UpdateLineup(ctx, req)
		return kitx.Response{Value: nil, Error: err}, nil
	}
}

func MakeDeleteLineupEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.DeleteLineupReq)
		err = s.DeleteLineup(ctx, req)
		return kitx.Response{Value: nil, Error: err}, nil
	}
}

func MakeNewMatchEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.NewMatchReq)
		res, err := s.NewMatch(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeGetMatchEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.GetMatchReq)
		res, err := s.GetMatch(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeJoinMatchEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.JoinMatchReq)
		res, err := s.JoinMatch(ctx, req)
		return kitx.Response{Value: res, Error: err}, nil
	}
}

func MakeSyncOperateEndpoint(s itfs.SignalingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*entity.SyncOperate)
		res, err := s.SyncOperate(ctx, req)
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
