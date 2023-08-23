package utils

import (
	"encoding/json"

	"github.com/spf13/cast"
)

const (
	StringifyMaxSize = 2048
)

func Stringify(i interface{}) string {
	str, err := cast.ToStringE(i)
	if err == nil {
		return str
	}

	b, err := json.Marshal(i)
	if err == nil {
		if len(b) <= StringifyMaxSize {
			return string(b)
		} else {
			return string(b[:StringifyMaxSize]) + "..."
		}
	}
	return ""
}
