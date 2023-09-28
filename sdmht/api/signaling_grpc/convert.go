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
		UnitsLocation:     in.GetUnitsLocation(),
		HandCards:         in.GetHandCards(),
		CardLibraries:     in.GetCardLibraries(),
		IsLibraryExpty:    in.GetIsLibraryExpty(),
		LibraryExptyHurt:  in.GetLibraryExptyHurt(),
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
		UnitsLocation:     in.UnitsLocation[:],
		HandCards:         in.HandCards[:],
		CardLibraries:     in.CardLibraries[:],
		IsLibraryExpty:    in.IsLibraryExpty,
		LibraryExptyHurt:  in.LibraryExptyHurt,
		DrawCardCountdown: in.DrawCardCountDown,
	}
}

func FromPBPlayer(in *pb.Player) (out *entity.Player) {
	if in == nil {
		return (*entity.Player)(nil)
	}
	out = &entity.Player{
		ID:    in.GetId(),
		Cost:  in.GetCost(),
		Scene: FromPBScene(in.GetScene()),
		Units: []*entity.Unit{},
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
		Id:    in.ID,
		Cost:  in.Cost,
		Scene: ToPBScene(in.Scene),
		Units: []*pb.Unit{},
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
		ID:        in.GetId(),
		SN:        in.GetSn(),
		Winner:    in.GetWinner(),
		WhoseTurn: in.GetWhoseTurn(),
		CurRound:  in.GetCurRound(),
		Players:   []*entity.Player{},
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
		Id:        in.ID,
		Sn:        in.SN,
		Winner:    in.Winner,
		WhoseTurn: in.WhoseTurn,
		CurRound:  in.CurRound,
		Players:   []*pb.Player{},
	}
	for _, player := range in.Players {
		out.Players = append(out.Players, ToPBPlayer(player))
	}
	return out
}

func FromPBUnit(in *pb.Unit) (out *entity.Unit) {
	if in == nil {
		return (*entity.Unit)(nil)
	}

	out = &entity.Unit{
		// Location:      in.GetLocation(),
		BaseAttribute: *FromPBBaseAttribute(in.GetBaseAttribute()),
		Health:        in.GetHealth(),
		Defend:        in.GetDefend(),
		Move:          in.GetMove(),
		AttackPrevent: in.GetAttackPrevent(),
		// CounterAttack: in.GetCounterAttack(),
		// Rebirth:       in.GetRebirth(),
		HurtInstead: in.GetHurtInstead(),
		NoMove:      in.GetNoMove(),
		NoAttack:    in.GetNoAttack(),
		NoCure:      in.GetNoCure(),
	}
	for _, v := range in.GetChangeAttack() {
		out.ChangeAttack = append(out.ChangeAttack, FromPBTempAttribute(v))
	}
	for _, v := range in.GetChangeMove() {
		out.ChangeMove = append(out.ChangeMove, FromPBTempAttribute(v))
	}
	for _, v := range in.GetHurt() {
		out.Hurt = append(out.Hurt, FromPBTempAttribute(v))
	}
	for _, v := range in.GetCure() {
		out.Cure = append(out.Cure, FromPBTempAttribute(v))
	}

	return out
}
func ToPBUnit(in *entity.Unit) (out *pb.Unit) {
	if in == nil {
		return (*pb.Unit)(nil)
	}
	out = &pb.Unit{
		// Location:      in.Location,
		BaseAttribute: ToPBBaseAttribute(&in.BaseAttribute),
		Health:        in.Health,
		Defend:        in.Defend,
		Move:          in.Move,
		AttackPrevent: in.AttackPrevent,
		// CounterAttack: in.CounterAttack,
		// Rebirth:       in.Rebirth,
		HurtInstead: in.HurtInstead,
		NoMove:      in.NoMove,
		NoAttack:    in.NoAttack,
		NoCure:      in.NoCure,
	}
	for _, v := range in.ChangeAttack {
		out.ChangeAttack = append(out.ChangeAttack, ToPBTempAttribute(v))
	}
	for _, v := range in.ChangeMove {
		out.ChangeMove = append(out.ChangeMove, ToPBTempAttribute(v))
	}
	for _, v := range in.Hurt {
		out.Hurt = append(out.Hurt, ToPBTempAttribute(v))
	}
	for _, v := range in.Cure {
		out.Cure = append(out.Cure, ToPBTempAttribute(v))
	}
	return out
}

func FromPBTempAttribute(in *pb.TempAttribute) (out *entity.TempAttribute) {
	if in == nil {
		return (*entity.TempAttribute)(nil)
	}
	return &entity.TempAttribute{
		Value:  in.GetValue(),
		Period: in.GetPeriod(),
	}
}
func ToPBTempAttribute(in *entity.TempAttribute) (out *pb.TempAttribute) {
	if in == nil {
		return (*pb.TempAttribute)(nil)
	}
	return &pb.TempAttribute{
		Value:  in.Value,
		Period: in.Period,
	}
}

func FromPBBaseAttribute(in *pb.BaseAttribute) (out *entity.BaseAttribute) {
	if in == nil {
		return (*entity.BaseAttribute)(nil)
	}
	return &entity.BaseAttribute{
		ID:           in.GetId(),
		Name:         in.GetName(),
		Rarity:       in.GetRarity(),
		Affiliate:    in.GetAffiliate(),
		Attack:       in.GetAttack(),
		MaxDefend:    in.GetMaxDefend(),
		MaxHealth:    in.GetMaxHealth(),
		MaxMove:      in.GetMaxMove(),
		BaseNoMove:   in.GetBaseNoMove(),
		BaseNoAttack: in.GetBaseNoAttack(),
		BaseNoCure:   in.GetBaseNoCure(),
		BaseNoEquip:  in.GetBaseNoEquip(),
	}
}
func ToPBBaseAttribute(in *entity.BaseAttribute) (out *pb.BaseAttribute) {
	if in == nil {
		return (*pb.BaseAttribute)(nil)
	}
	return &pb.BaseAttribute{
		Id:           in.ID,
		Name:         in.Name,
		Rarity:       in.Rarity,
		Affiliate:    in.Affiliate,
		Attack:       in.Attack,
		MaxDefend:    in.MaxDefend,
		MaxHealth:    in.MaxHealth,
		MaxMove:      in.MaxMove,
		BaseNoMove:   in.BaseNoMove,
		BaseNoAttack: in.BaseNoAttack,
		BaseNoCure:   in.BaseNoCure,
		BaseNoEquip:  in.BaseNoEquip,
	}
}
