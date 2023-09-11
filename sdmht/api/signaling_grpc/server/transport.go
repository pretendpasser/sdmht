package server

import (
	"context"

	"sdmht/lib/kitx"
	errpb "sdmht/lib/protobuf/types/error"
	grpc "sdmht/sdmht/api/signaling_grpc"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	entity "sdmht/sdmht/svc/entity"
)

func deLoginReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginReq)
	return &entity.LoginReq{
		WeChatID: req.GetWechatId(),
		UserName: req.GetUserName(),
	}, nil
}
func enLoginReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.LoginReply{}
	if r.Error != nil {
		res.Err = errpb.ToPbError(r.Error)
		return res, nil
	}
	rr := r.Value.(*entity.LoginRes)
	if rr == nil {
		return res, nil
	}
	res.AccountId = rr.AccountID
	return res, nil
}

func deNewLineupReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.NewLineupReq)
	return &entity.NewLineupReq{
		Lineup: *grpc.FromPBLineup(req.GetLineup()),
	}, nil
}

func deFindLineupReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.FindLineupReq)
	return &entity.FindLineupReq{
		AccountID: req.GetAccountId(),
	}, nil
}
func enFindLineupReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.FindLineupReply{}
	if r.Error != nil {
		res.Err = errpb.ToPbError(r.Error)
		return res, nil
	}
	rr := r.Value.(*entity.FindLineupRes)
	if rr == nil {
		return res, nil
	}
	res.Total = int32(rr.Total)
	for _, lineup := range rr.Lineups {
		res.Lineups = append(res.Lineups, grpc.ToPBLineup(lineup))
	}

	return res, nil
}

func deUpdateLineupReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateLineupReq)
	return &entity.UpdateLineupReq{
		Lineup: *grpc.FromPBLineup(req.GetLineup()),
	}, nil
}

func deDeleteLineupReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteLineupReq)
	return &entity.DeleteLineupReq{
		ID:        req.GetId(),
		AccountID: req.GetAccountId(),
	}, nil
}

func deNewMatchReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.NewMatchReq)
	return &entity.NewMatchReq{
		AccountID: req.GetAccountId(),
		Positions: req.GetPositions(),
		LineupID:  req.GetLineupId(),
	}, nil
}
func enNewMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.NewMatchReply{}
	if r.Error != nil {
		res.Err = errpb.ToPbError(r.Error)
		return res, nil
	}
	rr := r.Value.(*entity.NewMatchRes)
	if rr == nil {
		return res, nil
	}
	res.MatchId = rr.MatchID
	return res, nil
}

func enCommonReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	res := &pb.CommonReply{}
	if r.Error != nil {
		res.Err = errpb.ToPbError(r.Error)
		return res, nil
	}
	return res, nil
}

func deKeepAliveReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.KeepAliveReq)
	return &entity.KeepAliveReq{
		Operator: req.Operator,
	}, nil
}

func deLogoutReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LogoutReq)
	return &entity.LogoutReq{
		Operator: req.Operator,
		Reason:   req.Reason,
	}, nil
}
