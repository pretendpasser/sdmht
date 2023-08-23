package client

import (
	"context"

	"sdmht/lib"
	sdmht_entity "sdmht/sdmht/svc/entity"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func encodeDispatchEventToClientRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*sdmht_entity.ClientEvent)
	return &pb.ClientEventReq{
		UserId: req.UserID,
		Event: &pb.Event{
			Type:    req.Type,
			Content: string(req.Content),
			AtTime:  timestamppb.New(req.AtTime),
		},
	}, nil
}

func decodeDispatchEventToClientReply(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	gr := grpcResponse.(*pb.DispatchEventToClientReply)
	if gr.Err != nil {
		return nil, lib.NewError(int(gr.Err.Errno), gr.Err.Errmsg)
	}

	return sdmht_entity.DispatchEventToClientReply{
		UserID:    gr.ClientReply.UserId,
		OK:        gr.ClientReply.Ok,
		ClientErr: gr.ClientReply.Err,
	}, nil
}

func encodeKickClientRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(uint64)
	return &pb.KickClientReq{UserId: req}, nil
}

//func decodeKickClientReply(_ context.Context, request interface{}) (interface{}, error) {
//	return nil, nil
//}

func decodeCommonReply(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	gr := grpcResponse.(*pb.CommonReply)
	if gr.Err != nil {
		return nil, lib.NewError(int(gr.Err.Errno), gr.Err.Errmsg)
	}
	return nil, nil
}
