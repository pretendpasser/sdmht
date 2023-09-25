package grpc

import (
	"sdmht/account/api/grpc/pb"
	"sdmht/account/svc/entity"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertAccountFromPB(in *pb.Account) (out *entity.Account) {
	if in == nil {
		return nil
	}
	return &entity.Account{
		ID:          in.GetId(),
		WeChatID:    in.GetWechatId(),
		UserName:    in.GetUserName(),
		CreatedAt:   in.GetCreatedAt().AsTime(),
		LastLoginAt: in.GetLastLoginAt().AsTime(),
	}
}

func ConvertAccountToPB(in *entity.Account) (out *pb.Account) {
	if in == nil {
		return nil
	}

	return &pb.Account{
		Id:          in.ID,
		WechatId:    in.WeChatID,
		UserName:    in.UserName,
		CreatedAt:   timestamppb.New(in.CreatedAt),
		LastLoginAt: timestamppb.New(in.LastLoginAt),
	}
}
