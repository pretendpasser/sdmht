package entity

const (
	LinupSplitChar = ";"
	MaxCardLibrary = 20
)

type Lineup struct {
	ID              uint64   `json:"id" db:"id"`
	AccountID       uint64   `json:"account_id" db:"account_id"`
	Name            string   `json:"name" db:"name"`
	Enabled         bool     `json:"-" db:"enabled"`
	Units           []uint64 `json:"units" db:"-"`
	CardLibrarys    []uint64 `json:"card_librarys" db:"-"`
	UnitsStr        string   `json:"-" db:"units"`
	CardLibrarysStr string   `json:"-" db:"card_library"`
}

type LineupQuery struct {
	FilterByAccountID uint64
}
