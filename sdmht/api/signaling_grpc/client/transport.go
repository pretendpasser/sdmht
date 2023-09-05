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
		HandCard:          [10]int64{},
		CardLibrary:       [20]int64{},
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

func enLoginReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.LoginReq)
	return &pb.LoginReq{
		UserName: req.UserName,
		WechatId: req.WeChatID,
	}, nil
}

func deLoginReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.LoginReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return &entity.LoginRes{
		AccountID: r.GetAccountId(),
	}, nil
}

func enNewMatchReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.NewMatchReq)
	return &pb.NewMatchReq{
		Operator:   req.Operator,
		CardConfig: req.CardConfig,
	}, nil
}

func deNewMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.NewMatchReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return &entity.NewMatchRes{
		Player: fromPBPlayer(r.Player),
	}, nil
}

func enKeepAliveReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.KeepAliveReq)
	return &pb.KeepAliveReq{
		Operator: req.Operator,
	}, nil
}

func enLogoutReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.LogoutReq)
	return &pb.LogoutReq{
		Operator: req.Operator,
		Reason:   req.Reason,
	}, nil
}

func deCommonReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.CommonReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return nil, nil
}
