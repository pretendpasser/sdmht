package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ itfs.LineupRepo = (*lineupRepo)(nil)

type lineupRepo struct {
	db *sqlx.DB
}

func NewLineupRepo(db *sqlx.DB) *lineupRepo {
	return &lineupRepo{
		db: db,
	}
}

func (r *lineupRepo) Create(ctx context.Context, lineup *entity.Lineup) error {
	units, cardLibrarys := "", ""
	for _, unit := range lineup.Units {
		units = fmt.Sprintf("%s%d%s", units, unit, entity.LinupSplitChar)
	}
	for _, cardLibrary := range lineup.CardLibrarys {
		cardLibrarys = fmt.Sprintf("%s%d%s", cardLibrarys, cardLibrary, entity.LinupSplitChar)
	}

	if len(units) > 0 {
		lineup.UnitsStr = units[:len(units)-1]
	}
	if len(cardLibrarys) > 0 {
		lineup.CardLibrarysStr = cardLibrarys[:len(cardLibrarys)-1]
	}

	builder := sq.Insert(`lineup`).
		Columns("account_id", "name", "enabled",
			"units", "card_library").
		Values(lineup.AccountID, lineup.Name, lineup.Enabled,
			lineup.UnitsStr, lineup.CardLibrarysStr)

	sql, args, err := builder.ToSql()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}

func (r *lineupRepo) Get(ctx context.Context, accountID uint64, id uint64) (*entity.Lineup, error) {
	var res entity.Lineup
	err := r.db.GetContext(ctx, &res, "SELECT * FROM lineup WHERE id=? and account_id=?",
		id, accountID)
	if err != nil {
		return nil, err
	}

	units, cardLibrarys := strings.Split(res.UnitsStr, ";"), strings.Split(res.CardLibrarysStr, ";")
	for _, unit := range units {
		tmp, _ := strconv.ParseUint(unit, 10, 64)
		res.Units = append(res.Units, tmp)
	}
	for _, cardLibrary := range cardLibrarys {
		tmp, _ := strconv.ParseInt(cardLibrary, 10, 64)
		res.CardLibrarys = append(res.CardLibrarys, tmp)
	}

	return &res, nil
}

func (r *lineupRepo) Find(ctx context.Context, query *entity.LineupQuery) (total int, res []*entity.Lineup, err error) {
	builder := sq.Select().From(`lineup`)
	if query != nil {
		if query.FilterByAccountID != 0 {
			builder = builder.Where(sq.Eq{"id": query.FilterByAccountID})
		}
	}
	builder = builder.OrderBy("id desc")

	{
		sql, args, err := builder.Column("COUNT(*)").ToSql()
		if err != nil {
			return 0, nil, err
		}
		if err := r.db.GetContext(ctx, &total, sql, args...); err != nil {
			return 0, nil, err
		}
	}

	if total == 0 {
		return 0, []*entity.Lineup{}, nil
	}

	sql, args, err := builder.Column("*").ToSql()
	if err != nil {
		return 0, nil, err
	}
	if err := r.db.SelectContext(ctx, &res, sql, args...); err != nil {
		return 0, nil, err
	}

	for i := range res {
		units, cardLibrarys := strings.Split(res[i].UnitsStr, ";"), strings.Split(res[i].CardLibrarysStr, ";")
		for _, unit := range units {
			tmp, _ := strconv.ParseUint(unit, 10, 64)
			res[i].Units = append(res[i].Units, tmp)
		}
		for _, cardLibrary := range cardLibrarys {
			tmp, _ := strconv.ParseInt(cardLibrary, 10, 64)
			res[i].CardLibrarys = append(res[i].CardLibrarys, tmp)
		}
	}

	return total, res, nil
}

func (r *lineupRepo) Update(ctx context.Context, lineup *entity.Lineup) error {
	units, cardLibrarys := "", ""
	for _, unit := range lineup.Units {
		units = fmt.Sprintf("%s%d%s", units, unit, entity.LinupSplitChar)
	}
	for _, cardLibrary := range lineup.CardLibrarys {
		cardLibrarys = fmt.Sprintf("%s%d%s", cardLibrarys, cardLibrary, entity.LinupSplitChar)
	}

	if len(units) > 0 {
		lineup.UnitsStr = units[:len(units)-1]
	}
	if len(cardLibrarys) > 0 {
		lineup.CardLibrarysStr = cardLibrarys[:len(cardLibrarys)-1]
	}

	mset := make(map[string]interface{})
	mset["name"] = lineup.Name
	mset["units"] = lineup.UnitsStr
	mset["card_library"] = lineup.CardLibrarysStr
	mset["enabled"] = lineup.Enabled
	builder := sq.Update(`lineup`).Where(sq.Eq{"id": lineup.ID, "account_id": lineup.AccountID}).SetMap(mset)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}

func (r *lineupRepo) Delete(ctx context.Context, accountID uint64, id uint64) error {
	builder := sq.Delete(`lineup`).Where(sq.Eq{"id": id, "account_id": accountID})
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}
