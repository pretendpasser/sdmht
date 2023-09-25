package server

import (
	"context"
	"encoding/json"

	"sdmht/lib"
	errpb "sdmht/lib/protobuf/types/error"
	sdmht_entity "sdmht/sdmht/svc/entity"
	"sdmht/sdmht_conn/api"
	pb "sdmht/sdmht_conn/api/grpc/conn_pb"
)

func decodeDispatchEventToClientReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ClientEventReq)
	rr := sdmht_entity.ClientEvent{
		AccountID: req.GetAccountId(),
		Type:      req.GetEvent().GetType(),
		AtTime:    req.GetEvent().GetAtTime().AsTime(),
		Content:   json.RawMessage([]byte(req.GetEvent().GetContent())),
	}
	return rr, nil
}

func encodeDispatchEventToClientReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(api.Response)
	rsp := &pb.DispatchEventToClientReply{}

	if r.Error != nil {
		rsp.Err = toPbError(r.Error)
		return rsp, nil
	}

	rr := r.Value.(sdmht_entity.DispatchEventToClientReply)
	rsp.ClientReply = ConvertClientReplyToPB(rr)

	return rsp, nil
}

func decodeKickClientReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.KickClientReq)

	return req.UserId, nil
}

func encodeKickClientReply(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(api.Response)
	rsp := &pb.CommonReply{}

	if r.Error != nil {
		rsp.Err = toPbError(r.Error)
		return rsp, nil
	}

	return rsp, nil
}

func ConvertClientReplyToPB(reply sdmht_entity.DispatchEventToClientReply) *pb.ClientReply {
	return &pb.ClientReply{
		AccountId: reply.AccountID,
		Ok:        reply.OK,
		Err:       reply.ClientErr,
	}
}

func toPbError(err error) *errpb.Error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case lib.Error:
		return &errpb.Error{Errno: int32(e.Code), Errmsg: e.Message}
	case *lib.Error:
		return &errpb.Error{Errno: int32(e.Code), Errmsg: e.Message}
	default:
		return &errpb.Error{Errno: lib.ErrInternal, Errmsg: err.Error()}
	}
}
