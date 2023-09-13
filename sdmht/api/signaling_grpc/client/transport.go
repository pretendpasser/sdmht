package client

import (
	"context"

	"sdmht/lib"
	grpc "sdmht/sdmht/api/signaling_grpc"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/svc/entity"
)

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

func enNewLineupReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.NewLineupReq)
	return &pb.NewLineupReq{
		Lineup: grpc.ToPBLineup(&req.Lineup),
	}, nil
}

func enFindLineupReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.FindLineupReq)
	return &pb.FindLineupReq{
		AccountId: req.AccountID,
	}, nil
}
func deFindLineupReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.FindLineupReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}

	res := &entity.FindLineupRes{
		Total: int(r.Total),
	}
	for _, lineup := range r.Lineups {
		res.Lineups = append(res.Lineups, grpc.FromPBLineup(lineup))
	}
	return res, nil
}

func enUpdateLineupReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.UpdateLineupReq)
	return &pb.UpdateLineupReq{
		Lineup: grpc.ToPBLineup(&req.Lineup),
	}, nil
}

func enDeleteLineupReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.DeleteLineupReq)
	return &pb.DeleteLineupReq{
		Id:        req.ID,
		AccountId: req.AccountID,
	}, nil
}

func enNewMatchReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.NewMatchReq)
	return &pb.NewMatchReq{
		AccountId: req.AccountID,
		Positions: req.Positions,
		LineupId:  req.LineupID,
	}, nil
}
func deNewMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.NewMatchReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	return &entity.NewMatchRes{
		MatchID: r.MatchId,
	}, nil
}

func enGetMatchReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.GetMatchReq)
	return &pb.GetMatchReq{
		AccountId: req.AccountID,
	}, nil
}
func deGetMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.GetMatchReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	match := grpc.FromPBMatch(r.Match)
	return &entity.GetMatchRes{
		Match: *match,
	}, nil
}

func enJoinMatchReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*entity.JoinMatchReq)
	return &pb.JoinMatchReq{
		AccountId: req.AccountID,
		Positions: req.Positions,
		LineupId:  req.LineupID,
		MatchId:   req.MatchID,
	}, nil
}
func deJoinMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.JoinMatchReply)
	if r.GetErr() != nil {
		return nil, lib.NewError(int(r.Err.Errno), r.Err.Errmsg)
	}
	match := grpc.FromPBMatch(r.Match)
	return &entity.JoinMatchRes{
		Match: *match,
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
