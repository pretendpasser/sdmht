package entity

const (
	UserOffline string = "offline"
	UserOnline  string = "online"
)

type LoginReq struct {
	WeChatID string `json:"wechat_id"`
	UserName string `json:"username"`
}
type LoginRes struct {
	AccountID uint64 `json:"account_id"`
}

type NewLineupReq struct {
	Lineup
}

type GetLineupReq struct {
	ID        uint64 `json:"id"`
	AccountID uint64 `json:"account_id"`
}

type FindLineupReq struct {
	AccountID uint64 `json:"account_id"`
}
type FindLineupRes struct {
	Total   int       `json:"total"`
	Lineups []*Lineup `json:"lineups"`
}

type UpdateLineupReq struct {
	Lineup
}

type DeleteLineupReq struct {
	ID        uint64 `json:"id"`
	AccountID uint64 `json:"account_id"`
}

type NewMatchReq struct {
	AccountID uint64   `json:"account_id"`
	Positions []uint64 `json:"positions"`
	LineupID  uint64   `json:"lineup_id"`
}
type NewMatchRes struct {
	MatchID uint64 `json:"match_id"`
}

type JoinMatchReq struct {
	AccountID uint64  `json:"account_id"`
	MatchID   uint64  `json:"match_id"`
	Positions []int64 `json:"positions"`
}
type JoinMatchRes struct{}

type EndMatchReq struct{}
type EndMatchRes struct{}

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
