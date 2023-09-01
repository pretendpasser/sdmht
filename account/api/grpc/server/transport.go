package server

import (
	"context"

	"sdmht/account/api/grpc"
	"sdmht/account/api/grpc/pb"
	"sdmht/account/svc/entity"
	"sdmht/lib"
	"sdmht/lib/kitx"
	errorPB "sdmht/lib/protobuf/types/error"
)

func decodeRegisterReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterReq)
	return req, nil
}
func enRegisterRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.RegisterRes{}
	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			res.Err = &errorPB.Error{Errno: int32(err.Code), Errmsg: err.Message}
		} else {
			res.Err = &errorPB.Error{Errno: int32(lib.ErrInternal), Errmsg: r.Error.Error()}
		}
		return res, nil
	}
	return res, nil
}

func decodeLoginReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginReq)
	return req, nil
}
func enLoginRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.LoginRes{}
	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			res.Err = &errorPB.Error{Errno: int32(err.Code), Errmsg: err.Message}
		} else {
			res.Err = &errorPB.Error{Errno: int32(lib.ErrInternal), Errmsg: r.Error.Error()}
		}
		return res, nil
	}
	v := r.Value.(*entity.LoginRes)
	res.Token = v.Token
	res.Account = grpc.ConvertAccountToPB(v.Account)
	return res, nil
}

func decodeLogoutReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LogoutReq)
	return req.Token, nil
}
func enLogoutRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.LoginRes{}
	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			res.Err = &errorPB.Error{Errno: int32(err.Code), Errmsg: err.Message}
		} else {
			res.Err = &errorPB.Error{Errno: int32(lib.ErrInternal), Errmsg: r.Error.Error()}
		}
		return res, nil
	}
	return res, nil
}

func decodeAuthenticateReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AuthenticateReq)
	return req.Token, nil
}
func enAuthenticateRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.AuthenticateRes{}
	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			res.Err = &errorPB.Error{Errno: int32(err.Code), Errmsg: err.Message}
		} else {
			res.Err = &errorPB.Error{Errno: int32(lib.ErrInternal), Errmsg: r.Error.Error()}
		}
		return res, nil
	}
	v := r.Value.(*entity.Account)
	res.Account = grpc.ConvertAccountToPB(v)
	return res, nil
}

func decodeGetAccountReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAccountReq)
	return req, nil
}
func enGetAccountRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.GetAccountRes{}
	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			res.Err = &errorPB.Error{Errno: int32(err.Code), Errmsg: err.Message}
		} else {
			res.Err = &errorPB.Error{Errno: int32(lib.ErrInternal), Errmsg: r.Error.Error()}
		}
		return res, nil
	}
	v := r.Value.(*entity.Account)
	res.Account = grpc.ConvertAccountToPB(v)
	return res, nil
}
