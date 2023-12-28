package entity

const (
	UserOffline string = "offline"
	UserOnline  string = "online"

	OpEventAttack   string = "attack"       // from(unitid)    to(position)
	OpEventMove     string = "move"         // from(unitid)    to(position)
	OpEventCard     string = "card"         // from(cardid)    to(unitid)
	OpEventSkill    string = "active_skill" // from(skillname) to(empty)
	OpEventEndRound string = "end"          // from(playerid)  to(empty)
	// EventHurt       string = "hurt"         // from(unitid)    to(hurt num)
	// EventOpenSquare string = "open_square"
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
	AccountID uint64  `json:"account_id"`
	Positions []int64 `json:"positions"`
	LineupID  uint64  `json:"lineup_id"`
}
type NewMatchRes struct {
	MatchID uint64 `json:"match_id"`
}

type GetMatchReq struct {
	AccountID uint64 `json:"account_id"`
}
type GetMatchRes struct {
	Match
}

type JoinMatchReq struct {
	AccountID uint64  `json:"account_id"`
	MatchID   uint64  `json:"match_id"`
	Positions []int64 `json:"positions"`
	LineupID  uint64  `json:"lineup_id"`
}
type JoinMatchRes struct {
	Match
}

type SyncOperateReq struct {
	MatchID  uint64 `json:"match_id"`
	Operator int32  `json:"operator"` // player index: [0 1]
	Event    string `json:"event"`
	From     int64  `json:"from"`
	To       int64  `json:"to"`
}

type SyncOperateRes struct {
	Match
}

type KeepAliveReq struct {
	Operator uint64 `json:"operator"`
}

type LogoutReq struct {
	Operator uint64 `json:"operator"`
	Reason   string `json:"reason"`
}
