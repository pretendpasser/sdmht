package mw

import (
	"context"
	"net/http"
	"sdmht/lib"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/utils"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type contextKey string

var contextKeyAccountID = contextKey("account-id")

func AccountMW() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			accountID := r.Header.Get(runtime.MetadataHeaderPrefix + utils.AccountIDKey)
			if accountID == "" {
				kitx.ErrorResponse(lib.NewError(lib.ErrPermissionDenied, "not accountid"), w)
				return
			}
			ctx = SetAccountIDToContext(ctx, contextKeyAccountID, accountID)

			r = r.WithContext(ctx)

			log.S().Debugf("[AppAccess] accountId:%s", accountID)
			next.ServeHTTP(w, r)
		})
	}
}

func SetAccountIDToContext(ctx context.Context, key contextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetAccountIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(contextKeyAccountID).(string)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseUint(v, 10, 64)
	return id, (err == nil)
}
