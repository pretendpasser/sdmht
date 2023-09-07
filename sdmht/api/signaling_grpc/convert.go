package grpc

import (
	pb "sdmht/sdmht/api/signaling_grpc/signaling_pb"
	"sdmht/sdmht/svc/entity"
)

func FromPBLineup(in *pb.Lineup) (out *entity.Lineup) {
	if in == nil {
		return (*entity.Lineup)(nil)
	}
	out = &entity.Lineup{
		ID:           in.GetId(),
		AccountID:    in.GetAccountId(),
		Name:         in.GetName(),
		Enabled:      in.GetEnabled(),
		Units:        in.GetUnits(),
		CardLibrarys: in.GetCardLibrarys(),
	}
	return out
}
func ToPBLineup(in *entity.Lineup) (out *pb.Lineup) {
	if in == nil {
		return (*pb.Lineup)(nil)
	}
	res := &pb.Lineup{
		Id:        in.ID,
		AccountId: in.AccountID,
		Name:      in.Name,
		Enabled:   in.Enabled,
	}
	res.Units = append(res.Units, in.Units...)
	res.CardLibrarys = append(res.CardLibrarys, in.CardLibrarys...)

	return res
}

func FromPBScene(in *pb.Scene) (out *entity.Scene) {
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
func ToPBScene(in *entity.Scene) (out *pb.Scene) {
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

func FromPBPlayer(in *pb.Player) (out *entity.Player) {
	if in == nil {
		return (*entity.Player)(nil)
	}
	return &entity.Player{
		ID:    in.Id,
		Scene: FromPBScene(in.Scene),
	}
}
func ToPBPlayer(in *entity.Player) (out *pb.Player) {
	if in == nil {
		return (*pb.Player)(nil)
	}
	return &pb.Player{
		Id:    in.ID,
		Scene: ToPBScene(in.Scene),
	}
}
