package errors

import "github.com/pkg/errors"

func New(err string, args ...interface{}) error {
	return errors.Errorf(err, args...)
}

// Trace has tree forms:
// - Trace(err)
// - Trace(err, msg)
// - Trace(err, msg, args...)
func Trace(err error, args ...interface{}) error {
	if len(args) == 0 {
		return errors.WithStack(err)
	}
	msg := args[0].(string)
	return errors.Wrapf(err, msg, args[1:]...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
