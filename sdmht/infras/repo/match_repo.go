package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-redis/redis/v8"
)

var _ itfs.MatchRepo = (*matchRepo)(nil)

var (
	maxMatchTimeout = 300 * time.Second
)

type matchRepo struct {
	key          string
	match        map[uint64]*entity.Match
	matchTimeout map[uint64]*time.Timer
	rdb          *redis.Client
}

func NewMatchRepo(key string, rdb *redis.Client) *matchRepo {
	return &matchRepo{
		key:          key,
		match:        make(map[uint64]*entity.Match),
		matchTimeout: make(map[uint64]*time.Timer),
		rdb:          rdb,
	}
}

// func (r *matchRepo) CheckingTimeout() {
// 	for {
// 		for _, match := range r.matchTimeout {

// 		}
// 	}
// }

// ------------------------
//
//	Memory Cache
//
// ------------------------
func (r *matchRepo) New(match *entity.Match) error {
	if r.match[match.ID] != nil {
		return lib.NewError(lib.ErrInvalidArgument, "match exist")
	}
	r.match[match.ID] = match
	r.matchTimeout[match.ID] = time.NewTimer(maxMatchTimeout)
	go func(id uint64) {
		<-r.matchTimeout[id].C
		log.S().Infow("match has no updated in maxMatchTimeout", "matchid", id)
		r.matchTimeout[id].Stop()
		r.Delete(id)
	}(match.ID)
	return nil
}

func (r *matchRepo) Set(match *entity.Match) {
	r.match[match.ID] = match
	r.matchTimeout[match.ID].Reset(maxMatchTimeout)
}

func (r *matchRepo) Get(id uint64) (*entity.Match, error) {
	if r.match[id] == nil {
		return nil, lib.NewError(lib.ErrNotFound, "match not found")
	}
	return r.match[id], nil
}

func (r *matchRepo) Delete(id uint64) error {
	delete(r.match, id)
	delete(r.matchTimeout, id)
	return nil
}

// ------------------------
//
//	Redis
//
// ------------------------
func (r *matchRepo) RAdd(ctx context.Context, match *entity.Match) error {
	log.S().Infow("redis add match", "key", r.key, "id", match.ID)
	matchByte, err := json.Marshal(match)
	if err != nil {
		log.S().Errorw("redis add match marshal fail", "err", err)
		return err
	}
	_, err = r.rdb.HSet(ctx, r.key, match.ID, matchByte).Result()
	return err
}

func (r *matchRepo) RGet(ctx context.Context, id uint64) (*entity.Match, error) {
	ret, err := r.rdb.HGet(ctx, r.key, fmt.Sprintf("%d", id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.S().Errorw("redis get match not found", "err", err)
			return nil, lib.NewError(lib.ErrNotFound, "redis match not found")
		}
	}
	match := entity.Match{}
	err = json.Unmarshal(ret, &match)
	if err != nil {
		log.S().Errorw("redis get match unmarshal fail", "err", err)
		return nil, err
	}

	return &match, err
}

func (r *matchRepo) RUpdate(ctx context.Context, match *entity.Match) error {
	log.S().Infow("redis update match", "key", r.key, "id", match.ID)
	matchByte, err := json.Marshal(match)
	if err != nil {
		log.S().Errorw("redis update match marshal fail", "err", err)
		return err
	}
	_, err = r.rdb.HSet(ctx, r.key, match.ID, matchByte).Result()
	return err
}

func (r *matchRepo) RDelete(ctx context.Context, id uint64) error {
	_, err := r.rdb.HDel(ctx, r.key, fmt.Sprintf("%d", id)).Result()
	return err
}
