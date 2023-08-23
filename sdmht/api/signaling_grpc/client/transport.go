package client

import (
	"context"

	"sdmht/lib"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/svc/entity"
)

func fromPBScene(in *pb.Scene) (out *entity.Scene) {
	if in == nil {
		return (*entity.Scene)(nil)
	}
	out = &entity.Scene{
		Squares:           [16]int32{},
		HandCard:          [10]int32{},
		CardLibrary:       [20]int32{},
		DrawCardCountDown: in.DrawCardCountdown,
	}
	_ = copy(out.Squares[:], in.Squares)
	_ = copy(out.HandCard[:], in.HandCard)
	_ = copy(out.CardLibrary[:], in.CardLibrary)
	return out
}

func fromPBPlayer(in *pb.Player) (out *entity.Player) {
	if in == nil {
		return (*entity.Player)(nil)
	}
	return &entity.Player{
		ID:    in.Id,
		Scene: fromPBScene(in.Scene),
	}
}

func encodeNewMatchReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.NewMatchReq)
	return &pb.NewMatchReq{
		Operator:   req.Operator,
		CardConfig: req.CardConfig,
	}, nil
}

func decodeNewMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.NewMatchReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return &entity.NewMatchRsp{
		Player: fromPBPlayer(r.Player),
	}, nil
}

func encodeKeepAliveReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.KeepAliveReq)
	return &pb.KeepAliveReq{
		Operator: req.Operator,
	}, nil
}

func encodeLogoutReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.LogoutReq)
	return &pb.LogoutReq{
		Operator: req.Operator,
		Reason:   req.Reason,
	}, nil
}

func decodeCommonReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.CommonReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return nil, nil
}
