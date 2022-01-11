package domain

import (
	"errors"
	"fmt"
)

// nolint:errname
var Error = errors.New("domain error")

func NewError(msg string) error {
	return fmt.Errorf("%w: %s", Error, msg)
}
