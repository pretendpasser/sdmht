package lib

import (
	"fmt"
	"net/http"
)

const (
	ErrInvalidArgument   = 1
	ErrNotFound          = 2
	ErrAlreadyExists     = 3
	ErrPermissionDenied  = 4
	ErrResourceExhausted = 5
	ErrUnavailable       = 6
	ErrDataLoss          = 7
	ErrBusy              = 8
	ErrAborted           = 9
	ErrInternal          = 10
	ErrUnauthorized      = 11
	ErrUnsupported       = 12

	ErrSuccess = 200
)

var ErrorStrings = map[int]string{
	ErrInvalidArgument:   "invalid argument",
	ErrNotFound:          "not found",
	ErrAlreadyExists:     "aleady exists",
	ErrPermissionDenied:  "permission denied",
	ErrResourceExhausted: "resource exhausted",
	ErrUnavailable:       "unavailable",
	ErrDataLoss:          "data loss",
	ErrBusy:              "busy",
	ErrAborted:           "aborted",
	ErrInternal:          "internal",
	ErrUnauthorized:      "unauthorized",
	ErrUnsupported:       "unsupported",
	ErrSuccess:           "success",
}

type Error struct {
	Code    int    `json:"errno"`
	Message string `json:"errmsg"`
}

func NewError(c int, m string) Error {
	if m == "" {
		m = ErrorStrings[c]
	}
	return Error{c, m}
}

func (e Error) Error() string {
	return fmt.Sprintf("<Err(%d,%s)>", e.Code, e.Message)
}

func (e Error) HttpStatusCode() (statusCode int) {
	switch e.Code {
	case ErrInvalidArgument:
		statusCode = http.StatusBadRequest
	case ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case ErrNotFound:
		statusCode = http.StatusNotFound
	case ErrPermissionDenied:
		statusCode = http.StatusForbidden
	case ErrInternal:
		statusCode = http.StatusInternalServerError
	// case ErrAppIdKeyNotMatch, ErrInvalidJoinedToken, ErrPuberSub, ErrSuberPub:
	// statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}
