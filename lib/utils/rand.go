package utils

import (
	"math/rand"
	"sdmht/lib/log"
)

// support following type:
// []int, []int32. []int64, []uint, []uint32 ,[]uint64
func SliceRandom(s interface{}) interface{} {
	switch ss := s.(type) {
	case []int:
		res := make([]int, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	case []int32:
		res := make([]int32, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	case []int64:
		res := make([]int64, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	case []uint:
		res := make([]uint, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	case []uint32:
		res := make([]uint32, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	case []uint64:
		res := make([]uint64, len(ss))
		randList := rand.Perm(len(ss))
		for i, v := range randList {
			res[i] = ss[v]
		}
		return res
	}
	log.S().Errorw("invalid type", "s", s)
	return s
}
