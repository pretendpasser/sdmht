package entity

var (
	LinupSplitChar string = ";"
	MaxCardLibrary int    = 20
	MaxBaseUnitNum int    = 3
)

type Match struct {
	ID         uint64    `json:"id"`
	Winner     uint64    `json:"winner"`
	CurRoundID uint64    `json:"cur_round_id"`
	Players    []*Player `json:"players"`
}

type Player struct {
	ID     uint64  `json:"id"`
	MyTurn bool    `json:"my_turn"`
	Cost   int32   `json:"cost"`
	Scene  *Scene  `json:"scene"`
	Units  []*Unit `json:"units"`
}

type Lineup struct {
	ID              uint64  `json:"id" db:"id"`
	AccountID       uint64  `json:"account_id" db:"account_id"`
	Name            string  `json:"name" db:"name"`
	Enabled         bool    `json:"-" db:"enabled"`
	Units           []int64 `json:"units" db:"-"`
	CardLibrarys    []int64 `json:"card_librarys" db:"-"`
	UnitsStr        string  `json:"-" db:"units"`
	CardLibrarysStr string  `json:"-" db:"card_library"`
}

type LineupQuery struct {
	FilterByAccountID uint64
}
