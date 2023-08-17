package utils

import (
	"math/rand"
)

func SliceRandom(s *[]int) {
	randList := rand.Perm(len(*s))
	for i, v := range randList {
		randList[i] = (*s)[v]
	}
	*s = randList
}
