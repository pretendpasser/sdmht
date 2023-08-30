package utils

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	AccountIDKey     = "account-id"
	RootAccountIDKey = "root-account-id"
	AccountTypeKey   = "is-admin"
	NoCaptchaKey     = "login-no-captcha"
	ActionKey        = "http-method"
	ObjectKey        = "http-path"
	ForWardForIPKey  = "x-forwarded-for"
	RealIPKey        = "x-real-ip"
)

func GetGrpcMetadataFromContext(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	v, ok := md[key]
	if !ok || len(v) == 0 {
		return "", false
	}

	return v[0], true
}

func GetAccountIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := GetGrpcMetadataFromContext(ctx, AccountIDKey)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseUint(v, 10, 64)
	return id, (err == nil)
}

func GetRootAccountIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := GetGrpcMetadataFromContext(ctx, RootAccountIDKey)
	if !ok {
		return 0, false
	}
	// when current account's root account is itself, the value will be "0"
	if v == "0" {
		return GetAccountIDFromContext(ctx)
	}
	id, err := strconv.ParseUint(v, 10, 64)
	return id, (err == nil)
}

func GetIPFromContext(ctx context.Context) string {
	ip, ok := GetGrpcMetadataFromContext(ctx, ForWardForIPKey)
	if ok && ip != "" {
		return strings.TrimSpace(ip)
	}
	ip, ok = GetGrpcMetadataFromContext(ctx, RealIPKey)
	if ok && ip != "" {
		return strings.TrimSpace(ip)
	}

	return strings.TrimSpace(ip)
}

func CheckIsAdmin(ctx context.Context) bool {
	v, ok := GetGrpcMetadataFromContext(ctx, AccountTypeKey)
	if !ok {
		return false
	}
	isAdmin, _ := strconv.ParseBool(v)
	return isAdmin
}
