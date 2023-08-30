package http

import (
	"context"
	"encoding/json"
	"net/http"

	"sdmht/lib"
	"sdmht/sdmht_conn/api"
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	if err, ok := res.(error); ok {
		errorEncoder(ctx, err, w)
		return nil
	}
	if r, ok := res.(api.Response); ok && r.Error != nil {
		errorEncoder(ctx, r.Error, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(res.(api.Response).Value)
}

func errorEncoder(_ context.Context, failed error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	err, ok := failed.(lib.Error)
	if ok {
		switch err.Code {
		case lib.ErrNotFound:
			statusCode = http.StatusNotFound
		default:
		}
	} else {
		err = lib.Error{Code: lib.ErrInternal, Message: failed.Error()}
	}
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorWrapper{err})
}

type errorWrapper struct {
	Error error `json:"err"`
}
