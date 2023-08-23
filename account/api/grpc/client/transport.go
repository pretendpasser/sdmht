package client

import (
	"context"
	"sdmht/account/api/grpc"
	"sdmht/account/api/grpc/pb"
	"sdmht/lib"
)

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
