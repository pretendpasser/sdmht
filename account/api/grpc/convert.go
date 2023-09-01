package grpc

import (
	"sdmht/account/api/grpc/pb"
	"sdmht/account/svc/entity"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertAccountFromPB(account *pb.Account) *entity.Account {
	if account == nil {
		return nil
	}
	return &entity.Account{
		ID:       account.Id,
		UserName: account.UserName,
	}
}

func ConvertAccountToPB(account *entity.Account) *pb.Account {
	if account == nil {
		return nil
	}

	pbAccount := &pb.Account{
		Id:          account.ID,
		UserName:    account.UserName,
		CreatedAt:   timestamppb.New(account.CreatedAt),
		LastLoginAt: timestamppb.New(account.LastLoginAt),
	}

	return pbAccount
}
