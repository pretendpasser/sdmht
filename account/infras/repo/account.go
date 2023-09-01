package repo

import (
	"context"

	"sdmht/account/svc/entity"
	itfs "sdmht/account/svc/interfaces"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ itfs.AccountRepo = (*accountRepo)(nil)

type accountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) Add(ctx context.Context, account *entity.Account) error {
	builder := sq.Insert(`account`).
		Columns("user_name", "wechat_id").
		Values(account.UserName, account.WeChatID)
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) Get(ctx context.Context, id uint64) (*entity.Account, error) {
	builder := sq.Select().From(`account`).Where(sq.Eq{"id": id})
	sql, args, err := builder.Column("*").ToSql()
	if err != nil {
		return nil, err
	}
	var account entity.Account
	if err := r.db.GetContext(ctx, &account, sql, args...); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) GetByWechatID(ctx context.Context, wechatID string) (*entity.Account, error) {
	builder := sq.Select().From(`account`).Where(sq.Eq{"wechat_id": wechatID})
	sql, args, err := builder.Column("*").ToSql()
	if err != nil {
		return nil, err
	}
	var account entity.Account
	if err := r.db.GetContext(ctx, &account, sql, args...); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) Update(ctx context.Context, account *entity.Account) error {
	mset := make(map[string]interface{})
	mset["user_name"] = account.UserName
	mset["wechat_id"] = account.WeChatID
	mset["last_login_at"] = account.LastLoginAt

	builder := sq.Update(`account`).Where(sq.Eq{"id": account.ID}).SetMap(mset)
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) Delete(ctx context.Context, id uint64) error {
	builder := sq.Delete(`account`).Where(sq.Eq{"id": id})
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	var account entity.Account
	if err := r.db.GetContext(ctx, &account, sql, args...); err != nil {
		return err
	}

	return nil
}
