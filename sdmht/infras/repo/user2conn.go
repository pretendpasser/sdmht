package repo

import (
	"context"
	"fmt"

	"sdmht/lib"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-redis/redis/v8"
)

var _ itfs.User2ConnRepo = (*user2ConnRepo)(nil)

type user2ConnRepo struct {
	key string
	rdb *redis.Client
}

func NewUser2ConnRepo(key string, rdb *redis.Client) itfs.User2ConnRepo {
	return &user2ConnRepo{key, rdb}
}

func (r *user2ConnRepo) Find(ctx context.Context, neids []uint64) (map[uint64]string, error) {
	fields := make([]string, 0, len(neids))
	for _, neid := range neids {
		fields = append(fields, fmt.Sprintf("%d", neid))
	}

	vals, err := r.rdb.HMGet(ctx, r.key, fields...).Result()
	if err != nil {
		return nil, err
	}

	res := map[uint64]string{}
	for i, val := range vals { // FIXME field is not exists?
		s, ok := val.(string)
		if !ok {
			continue
		}
		res[neids[i]] = s
	}

	return res, nil
}

func (r *user2ConnRepo) Get(ctx context.Context, id uint64) (string, error) {
	field := fmt.Sprintf("%d", id)
	ret, err := r.rdb.HGet(ctx, r.key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", lib.NewError(lib.ErrNotFound, "")
		}
		return "", err
	}

	return ret, nil
}

func (r *user2ConnRepo) FindConnNames(ctx context.Context) ([]string, error) {
	return r.rdb.HVals(ctx, r.key).Result()
}

func (r *user2ConnRepo) Add(ctx context.Context, id uint64, wechat_id string) error {
	fmt.Println(r.key, id)
	_, err := r.rdb.HSet(ctx, r.key, id, wechat_id).Result()
	return err
}

func (r *user2ConnRepo) Delete(ctx context.Context, id uint64) error {
	_, err := r.rdb.HDel(ctx, r.key, fmt.Sprintf("%d", id)).Result()
	return err
}
