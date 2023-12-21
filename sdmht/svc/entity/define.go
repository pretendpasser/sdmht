package entity

var (
	LinupSplitChar string = ";"
	MaxCardLibrary int    = 20
	MaxBaseUnitNum int    = 3
)

type Match struct {
	ID        uint64   `json:"id"`
	SN        int64    `json:"sn"`
	Winner    uint64   `json:"winner"`
	WhoseTurn int32    `json:"whose_turn"` // player index: [0 1]
	CurRound  uint64   `json:"cur_round"`
	Scenes    []*Scene `json:"scenes"`
}

func (m *Match) GetOtherPlayer() int32 {
	if m.WhoseTurn == 0 {
		return 1
	}
	return 0
}

type Lineup struct {
	ID              uint64  `json:"id" db:"id"`
	AccountID       uint64  `json:"account_id" db:"account_id"`
	Name            string  `json:"name" db:"name"`
	Enabled         bool    `json:"enabled" db:"enabled"`
	Units           []int64 `json:"units" db:"-"`
	CardLibrarys    []int64 `json:"card_librarys" db:"-"`
	UnitsStr        string  `json:"-" db:"units"`
	CardLibrarysStr string  `json:"-" db:"card_library"`
}

type LineupQuery struct {
	FilterByAccountID uint64
}
