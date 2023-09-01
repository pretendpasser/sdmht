package entity

const (
	UserOffline string = "offline"
	UserOnline  string = "online"
)

type LoginReq struct{}

type NewMatchReq struct {
	Operator   uint64 `json:"operator"`
	CardConfig uint64 `json:"card_config"`
}

type NewMatchRsp struct {
	Player *Player `json:"player"`
}

// type JoinMatchReq struct {
// 	Operator   uint64 `json:"operator"`
// 	MatchID    uint64 `json:"match_id"`
// 	CardConfig uint64 `json:"card_config"`
// }

// type JoinMatchRes struct {
// 	Operator   uint64 `json:"operator"`
// 	CardConfig uint64 `json:"card_config"`
// }

type KeepAliveReq struct {
	Operator uint64 `json:"operator"`
}

type LogoutReq struct {
	Operator uint64 `json:"operator"`
	Reason   string `json:"reason"`
}
