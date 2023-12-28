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
		PlayerID:          in.GetPlayerId(),
		MasterID:          in.GetMasterId(),
		Units:             make(map[int64]*entity.Unit),
		Squares:           in.GetSquares(),
		UnitsLocation:     in.GetUnitsLocation(),
		HandCards:         in.GetHandCards(),
		CardLibraries:     in.GetCardLibraries(),
		IsLibraryExpty:    in.GetIsLibraryExpty(),
		RetainerAliveNum:  in.GetRetainerAliveNum(),
		LibraryExptyHurt:  in.GetLibraryExptyHurt(),
		DrawCardCountDown: in.GetDrawCardCountdown(),
		Cost:              in.GetCost(),
		LastMoveUnitID:    in.GetLastMoveUnitId(),
	}
	for _, unitPB := range in.GetUnits() {
		unit := FromPBUnit(unitPB)
		out.Units[unit.ID] = unit
	}
	return out
}
func ToPBScene(in *entity.Scene) (out *pb.Scene) {
	if in == nil {
		return (*pb.Scene)(nil)
	}
	out = &pb.Scene{
		PlayerId:          in.PlayerID,
		MasterId:          in.MasterID,
		Squares:           in.Squares[:],
		UnitsLocation:     in.UnitsLocation[:],
		HandCards:         in.HandCards[:],
		CardLibraries:     in.CardLibraries[:],
		IsLibraryExpty:    in.IsLibraryExpty,
		RetainerAliveNum:  in.RetainerAliveNum,
		LibraryExptyHurt:  in.LibraryExptyHurt,
		DrawCardCountdown: in.DrawCardCountDown,
		Cost:              in.Cost,
		LastMoveUnitId:    in.LastMoveUnitID,
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
		Scenes:    []*entity.Scene{},
	}
	for _, scene := range in.GetScenes() {
		out.Scenes = append(out.Scenes, FromPBScene(scene))
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
		Scenes:    []*pb.Scene{},
	}
	for _, scene := range in.Scenes {
		out.Scenes = append(out.Scenes, ToPBScene(scene))
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
		Attack:        in.GetAttack(),
		Move:          in.GetMove(),
		IsMoving:      in.GetIsMoving(),
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
		Attack:        in.Attack,
		Move:          in.Move,
		IsMoving:      in.IsMoving,
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
		BaseAttack:   in.GetBaseAttack(),
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
		BaseAttack:   in.BaseAttack,
		MaxDefend:    in.MaxDefend,
		MaxHealth:    in.MaxHealth,
		MaxMove:      in.MaxMove,
		BaseNoMove:   in.BaseNoMove,
		BaseNoAttack: in.BaseNoAttack,
		BaseNoCure:   in.BaseNoCure,
		BaseNoEquip:  in.BaseNoEquip,
	}
}
