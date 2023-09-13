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
		Squares:           in.GetSquares(),
		HandCards:         in.GetHandCards(),
		CardLibraries:     in.GetHandCards(),
		IsLibraryExpty:    in.GetIsLibraryExpty(),
		DrawCardCountDown: in.GetDrawCardCountdown(),
	}
	return out
}
func ToPBScene(in *entity.Scene) (out *pb.Scene) {
	if in == nil {
		return (*pb.Scene)(nil)
	}
	return &pb.Scene{
		Squares:           in.Squares[:],
		HandCards:         in.HandCards[:],
		CardLibraries:     in.CardLibraries[:],
		IsLibraryExpty:    in.IsLibraryExpty,
		DrawCardCountdown: in.DrawCardCountDown,
	}
}

func FromPBPlayer(in *pb.Player) (out *entity.Player) {
	if in == nil {
		return (*entity.Player)(nil)
	}
	out = &entity.Player{
		ID:     in.GetId(),
		MyTurn: in.GetMyTurn(),
		Cost:   in.GetCost(),
		Scene:  FromPBScene(in.GetScene()),
		Units:  []*entity.Unit{},
	}
	for _, unit := range in.GetUnits() {
		out.Units = append(out.Units, FromPBUnit(unit))
	}
	return out
}
func ToPBPlayer(in *entity.Player) (out *pb.Player) {
	if in == nil {
		return (*pb.Player)(nil)
	}
	out = &pb.Player{
		Id:     in.ID,
		MyTurn: in.MyTurn,
		Cost:   in.Cost,
		Scene:  ToPBScene(in.Scene),
		Units:  []*pb.Unit{},
	}
	for _, unit := range in.Units {
		out.Units = append(out.Units, ToPBUnit(unit))
	}
	return out
}

func FromPBMatch(in *pb.Match) (out *entity.Match) {
	if in == nil {
		return (*entity.Match)(nil)
	}
	out = &entity.Match{
		ID:         in.GetId(),
		Winner:     in.GetWinner(),
		CurRoundID: in.GetCurRoundId(),
		Players:    []*entity.Player{},
	}
	for _, player := range in.GetPlayers() {
		out.Players = append(out.Players, FromPBPlayer(player))
	}
	return out
}
func ToPBMatch(in *entity.Match) (out *pb.Match) {
	if in == nil {
		return (*pb.Match)(nil)
	}
	out = &pb.Match{
		Id:         in.ID,
		Winner:     in.Winner,
		CurRoundId: in.CurRoundID,
		Players:    []*pb.Player{},
	}
	for _, player := range in.Players {
		out.Players = append(out.Players, ToPBPlayer(player))
	}
	return out
}

func FromPBUnit(in *pb.Unit) (out *entity.Unit) { return }
func ToPBUnit(in *entity.Unit) (out *pb.Unit)   { return }
