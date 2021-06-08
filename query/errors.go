package query

import (
	"errors"
	"strings"
)

var ErrNotFound = errors.New("not found")

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not found")
}
