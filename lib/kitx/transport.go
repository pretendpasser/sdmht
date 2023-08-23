package kitx

import (
	"context"
	"encoding/json"
	"net/http"

	"sdmht/lib"

	"github.com/go-kit/kit/sd/lb"
)

// CodecEmpty for gRpc transport
func CodecEmpty(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

// ErrorResponse for http transport
func ErrorResponse(failed error, w http.ResponseWriter) {
	var (
		err      lib.Error
		isLibErr bool
	)

	if retryerr, retry := failed.(lb.RetryError); retry {
		err, isLibErr = retryerr.Final.(lib.Error)
	} else {
		err, isLibErr = failed.(lib.Error)
	}

	if !isLibErr {
		err = lib.Error{Code: lib.ErrInternal, Message: failed.Error()}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.HttpStatusCode())
	_ = json.NewEncoder(w).Encode(err)
}
