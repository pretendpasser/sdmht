package client

import (
	"context"
	"sdmht/account/api/grpc"
	"sdmht/account/api/grpc/pb"
	"sdmht/account/svc/entity"
	"sdmht/lib"
)

func enRegisterReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.RegisterReq)
	return &pb.LoginReq{
		WechatId: req.WechatID,
		UserName: req.UserName,
	}, nil
}

func deRegisterRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.RegisterRes)
	if err := res.GetErr(); err != nil {
		return nil, lib.NewError(int(err.Errno), err.Errmsg)
	}
	return nil, nil
}

func enLoginReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.LoginReq)
	return &pb.LoginReq{
		WechatId: req.WechatID,
		UserName: req.UserName,
	}, nil
}

func deLoginRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.LoginRes)
	if err := res.GetErr(); err != nil {
		return nil, lib.NewError(int(err.Errno), err.Errmsg)
	}
	return &entity.LoginRes{Token: res.GetToken()}, nil
}

func enLogoutReq(_ context.Context, request interface{}) (interface{}, error) {
	token := request.(string)
	return &pb.LogoutReq{
		Token: token,
	}, nil
}

func deLogoutRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.LogoutRes)
	if err := res.GetErr(); err != nil {
		return nil, lib.NewError(int(err.Errno), err.Errmsg)
	}
	return nil, nil
}

func enAuthenticateReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(string)
	return &pb.AuthenticateReq{Token: req}, nil
}

func deAuthenticateRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.AuthenticateRes)
	if err := res.GetErr(); err != nil {
		return nil, lib.NewError(int(err.Errno), err.Errmsg)
	}
	return grpc.ConvertAccountFromPB(res.GetAccount()), nil
}

func enGetAccountReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(uint64)
	return &pb.GetAccountReq{
		Id: req,
	}, nil
}

func deGetAccountRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.GetAccountRes)
	if err := res.GetErr(); err != nil {
		return nil, lib.NewError(int(err.Errno), err.Errmsg)
	}
	return grpc.ConvertAccountFromPB(res.GetAccount()), nil
}
