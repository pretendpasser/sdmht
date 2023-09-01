package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib"

	"github.com/go-redis/redis/v8"
)

var _ itfs.TokenRepo = (*tokenRepo)(nil)

type tokenRepo struct {
	rdb    *redis.Client
	prefix string
}

func NewTokenRepo(rdb *redis.Client, prefix string) *tokenRepo {
	return &tokenRepo{
		rdb:    rdb,
		prefix: prefix,
	}
}

func (r *tokenRepo) Add(ctx context.Context, token string, accountID uint64, ttl time.Duration) error {
	key := r.genKey(token)
	_, err := r.rdb.Set(ctx, key, accountID, ttl).Result()
	return err
}

func (r *tokenRepo) Get(ctx context.Context, token string, ttl time.Duration) (uint64, error) {
	key := r.genKey(token)
	v, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, lib.NewError(lib.ErrUnauthorized, "no token found")
		}
		return 0, err
	}
	_, _ = r.rdb.Expire(ctx, key, ttl).Result()
	return strconv.ParseUint(v, 10, 64)
}

func (r *tokenRepo) Delete(ctx context.Context, token string) error {
	key := r.genKey(token)
	_, err := r.rdb.Del(ctx, key).Result()
	return err
}

func (r *tokenRepo) genKey(key string) string {
	return fmt.Sprintf("%s%s", r.prefix, key)
}
