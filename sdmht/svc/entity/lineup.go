package entity

const (
	LinupSplitChar = ";"
	MaxCardLibrary = 20
)

type Lineup struct {
	ID              uint64   `db:"id"`
	AccountID       uint64   `db:"account_id"`
	Name            string   `db:"name"`
	Enabled         bool     `db:"enabled"`
	Units           []uint64 `db:"-"`
	CardLibrarys    []uint64 `db:"-"`
	UnitsStr        string   `db:"units"`
	CardLibrarysStr string   `db:"card_library"`
}

type LineupQuery struct {
	FilterByAccountID uint64
}
