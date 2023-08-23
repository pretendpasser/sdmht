package signaling_grpc

import (
	"context"

	"sdmht/lib/kitx"
	errpb "sdmht/lib/protobuf/types/error"
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	entity "sdmht/sdmht/svc/entity"
)

func toPBScene(in *entity.Scene) (out *pb.Scene) {
	if in == nil {
		return (*pb.Scene)(nil)
	}
	return &pb.Scene{
		Squares:           in.Squares[:],
		HandCard:          in.HandCard[:],
		CardLibrary:       in.CardLibrary[:],
		DrawCardCountdown: in.DrawCardCountDown,
	}
}

func toPBPlayer(in *entity.Player) (out *pb.Player) {
	if in == nil {
		return (*pb.Player)(nil)
	}
	return &pb.Player{
		Id:    in.ID,
		Scene: toPBScene(in.Scene),
	}
}

func decodeNewMatchReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.NewMatchReq)
	return &entity.NewMatchReq{
		Operator:   req.Operator,
		CardConfig: req.CardConfig,
	}, nil
}

func encodeNewMatchReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	rsp := &pb.NewMatchReply{}
	if r.Error != nil {
		rsp.Err = errpb.ToPbError(r.Error)
		return rsp, nil
	}
	rr := r.Value.(*entity.NewMatchRsp)
	if rr == nil {
		return rsp, nil
	}
	rsp.Player = toPBPlayer(rr.Player)
	return rsp, nil
}

func encodeCommonReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(kitx.Response)
	rsp := &pb.CommonReply{}
	if r.Error != nil {
		rsp.Err = errpb.ToPbError(r.Error)
		return rsp, nil
	}
	return rsp, nil
}

func decodeKeepAliveReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.KeepAliveReq)
	return &entity.KeepAliveReq{
		Operator: req.Operator,
	}, nil
}

func decodeLogoutReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LogoutReq)
	return &entity.LogoutReq{
		Operator: req.Operator,
		Reason:   req.Reason,
	}, nil
}
