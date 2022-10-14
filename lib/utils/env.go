package utils

import "os"

func GetEnvDefault(key, defalut string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defalut
	}
	return val
}
