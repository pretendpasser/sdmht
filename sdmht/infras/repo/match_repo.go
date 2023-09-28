package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/go-redis/redis/v8"
)

var _ itfs.MatchRepo = (*matchRepo)(nil)

var (
	maxMatchTimeout = 300 * time.Minute
)

type matchRepo struct {
	mux          sync.RWMutex
	key          string
	account      map[uint64]uint64
	match        map[uint64]*entity.Match
	matchTimeout map[uint64]*time.Timer
	rdb          *redis.Client
}

func NewMatchRepo(key string, rdb *redis.Client) *matchRepo {
	return &matchRepo{
		mux:          sync.RWMutex{},
		key:          key,
		account:      make(map[uint64]uint64),
		match:        make(map[uint64]*entity.Match),
		matchTimeout: make(map[uint64]*time.Timer),
		rdb:          rdb,
	}
}

// ------------------------
//
//	Memory Cache
//
// ------------------------
func (r *matchRepo) SetByAccount(accountID uint64, matchID uint64) {
	r.account[accountID] = matchID
}

func (r *matchRepo) GetByAccount(accountID uint64) (matchID uint64) {
	return r.account[accountID]
}

func (r *matchRepo) New(match *entity.Match) error {
	if r.match[match.ID] != nil {
		return lib.NewError(lib.ErrInvalidArgument, "match exist")
	}
	if len(match.Players) == 0 {
		return lib.NewError(lib.ErrInvalidArgument, "no player")
	}
	r.match[match.ID] = match
	r.matchTimeout[match.ID] = time.NewTimer(maxMatchTimeout)
	go func(id uint64) {
		<-r.matchTimeout[id].C
		log.S().Infow("match has no updated in maxMatchTimeout", "matchid", id)
		r.Delete(id)
	}(match.ID)
	return nil
}

func (r *matchRepo) Join(match *entity.Match) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if len(r.match[match.ID].Players) >= 2 {
		return lib.NewError(lib.ErrInternal, "match is already full")
	} else if len(r.match[match.ID].Players) == 0 {
		return lib.NewError(lib.ErrInternal, "match is invalid")
	}
	r.match[match.ID] = match
	r.matchTimeout[match.ID].Reset(maxMatchTimeout)
	return nil
}

func (r *matchRepo) Set(match *entity.Match) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.match[match.ID] = match
	r.matchTimeout[match.ID].Reset(maxMatchTimeout)
}

func (r *matchRepo) Get(id uint64) (entity.Match, error) {
	if r.match[id] == nil {
		return entity.Match{}, lib.NewError(lib.ErrNotFound, "match not found")
	}
	return *r.match[id], nil
}

func (r *matchRepo) Delete(id uint64) {
	r.mux.Lock()
	defer r.mux.Unlock()
	if r.matchTimeout[id] != nil {
		r.matchTimeout[id].Stop()
		delete(r.matchTimeout, id)
	}
	if r.match[id] != nil {
		for _, player := range r.match[id].Players {
			delete(r.account, player.ID)
		}
		delete(r.match, id)
	}
}

// ------------------------
//
//	Redis
//
// ------------------------
func (r *matchRepo) genAccountKey(accountID uint64) string {
	return fmt.Sprintf("%s:%d", r.key, accountID)
}

func (r *matchRepo) RSet(ctx context.Context, match *entity.Match) error {
	for _, player := range match.Players {
		key := r.genAccountKey(player.ID)
		_, err := r.rdb.Set(ctx, key, match.ID, maxMatchTimeout).Result()
		if err != nil {
			log.S().Errorw("redis set account-match not found", "err", err)
			return err
		}
	}
	return nil
}

func (r *matchRepo) RGet(ctx context.Context, accountID uint64) (uint64, error) {
	key := r.genAccountKey(accountID)
	ret, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.S().Errorw("redis get account-match not found", "err", err)
			return 0, lib.NewError(lib.ErrNotFound, "redis account-match not found")
		}
	}
	matchID, err := strconv.ParseUint(ret, 10, 64)
	if err != nil {
		log.S().Errorw("redis get account-match parseuint fail", "err", err)
		return 0, err
	}
	return matchID, nil
}

func (r *matchRepo) RDelete(ctx context.Context, accountID uint64) error {
	key := r.genAccountKey(accountID)
	_, err := r.rdb.Del(ctx, key).Result()
	return err
}

func (r *matchRepo) RHSet(ctx context.Context, match *entity.Match) error {
	log.S().Infow("redis hset match", "key", r.key, "id", match.ID)
	matchByte, err := json.Marshal(match)
	if err != nil {
		log.S().Errorw("redis hadd match marshal fail", "err", err)
		return err
	}
	_, err = r.rdb.HSet(ctx, r.key, match.ID, matchByte).Result()
	return err
}

func (r *matchRepo) RHGet(ctx context.Context, id uint64) (*entity.Match, error) {
	ret, err := r.rdb.HGet(ctx, r.key, fmt.Sprintf("%d", id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.S().Errorw("redis hget match not found", "err", err)
			return nil, lib.NewError(lib.ErrNotFound, "redis match not found")
		}
	}
	match := entity.Match{}
	err = json.Unmarshal(ret, &match)
	if err != nil {
		log.S().Errorw("redis hget match unmarshal fail", "err", err)
		return nil, err
	}

	return &match, err
}

func (r *matchRepo) RHDelete(ctx context.Context, id uint64) error {
	_, err := r.rdb.HDel(ctx, r.key, fmt.Sprintf("%d", id)).Result()
	return err
}
