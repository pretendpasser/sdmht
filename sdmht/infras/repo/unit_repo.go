package repo

import (
	"context"
	"fmt"
	"sdmht/lib"
	"strings"

	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ itfs.UnitRepo = (*unitRepo)(nil)

type unitRepo struct {
	db        *sqlx.DB
	skillList entity.SkillList
}

func NewUnitRepo(db *sqlx.DB, skillList entity.SkillList) *unitRepo {
	return &unitRepo{
		db:        db,
		skillList: skillList,
	}
}

func (r *unitRepo) Get(ctx context.Context, ids []int64) ([]*entity.Unit, error) {
	if len(ids) == 0 {
		return nil, lib.NewError(lib.ErrInvalidArgument, "no gods")
	}
	filter := ""
	for _, id := range ids {
		filter = fmt.Sprintf(`%s%d,`, filter, id)
	}
	filter = filter[:len(filter)-1]

	cmd := fmt.Sprintf(`SELECT * FROM god WHERE id in (%s)`, filter)

	var units []*entity.Unit
	err := r.db.SelectContext(ctx, &units, cmd)
	if err != nil {
		return nil, err
	}

	for i, unit := range units {
		units[i].Skills = make(map[string]*entity.Skill)
		skillNames := strings.Split(unit.SkillName, ";")
		for _, skillName := range skillNames {
			units[i].Skills[skillName] = &entity.Skill{
				Desc:    r.skillList.Get(skillName),
				Handler: r.skillList.Checking(skillName),
			}
		}
	}

	return units, nil
}

func (r *unitRepo) Find(ctx context.Context, query *entity.UnitQuery) (total int, units []*entity.Unit, err error) {
	builder := sq.Select().From(`god`)
	if query != nil {
		if query.Pagination != nil {
			limit, offset := query.Pagination.LimitOffset()
			builder = builder.Limit(limit).Offset(offset)
		}
		if query.ExcludeMainDeity {
			builder = builder.Where(sq.NotEq{"rarity": 0})
		}
		if query.FilterByRarity != 0 {
			builder = builder.Where(sq.Eq{"rarity": query.FilterByRarity})
		}
		if query.FilterByAffiliate != 0 {
			builder = builder.Where(sq.Eq{"affiliate": query.FilterByAffiliate})
		}
	}
	builder = builder.OrderBy("id asc")

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
		return 0, []*entity.Unit{}, nil
	}

	sql, args, err := builder.Column("*").ToSql()
	if err != nil {
		return 0, nil, err
	}
	if err := r.db.SelectContext(ctx, &units, sql, args...); err != nil {
		return 0, nil, err
	}

	for i, unit := range units {
		units[i].Skills = make(map[string]*entity.Skill)
		skillNames := strings.Split(unit.SkillName, ";")
		for _, skillName := range skillNames {
			units[i].Skills[skillName] = &entity.Skill{
				Desc:    r.skillList.Get(skillName),
				Handler: r.skillList.Checking(skillName),
			}
		}
	}

	return total, units, nil
}
