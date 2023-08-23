package error

import "sdmht/lib"

func ToPbError(err error) *Error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case lib.Error:
		return &Error{Errno: int32(e.Code), Errmsg: e.Message}
	case *lib.Error:
		return &Error{Errno: int32(e.Code), Errmsg: e.Message}
	default:
		return &Error{Errno: lib.ErrInternal, Errmsg: err.Error()}
	}
}
