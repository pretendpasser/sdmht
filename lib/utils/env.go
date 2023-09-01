package utils

import "os"

var gJWTSecret string = "example_jwt_secret"

func init() {
	gJWTSecret = GetEnvDefault("JWT_SECRET", gJWTSecret)
}

func GetJWTSecret() string {
	return gJWTSecret
}

func GetEnvDefault(key, defalut string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defalut
	}
	return val
}
