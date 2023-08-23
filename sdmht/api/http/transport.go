package http

import (
	"context"
	"encoding/json"
	"net/http"

	"sdmht/lib"
)

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
