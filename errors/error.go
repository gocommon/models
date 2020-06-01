package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error Error
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// New New
func New(msg string, args ...interface{}) error {
	return &Error{fmt.Sprintf(msg, args...)}
}

// Wrap Wrap
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*Error); !ok {
		err = &Error{err.Error()}
	}
	return errors.Wrap(err, msg)
}

// IsModelErr  IsModelErr
func IsModelErr(err error) bool {
	_, ok := err.(*Error)
	return ok
}
