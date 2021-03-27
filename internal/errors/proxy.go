package errors

import (
	"github.com/pkg/errors"
)

// PkgNew is as the proxy for github.com/pkg/errors.New func.
// func New(message string) error {
// 	return errors.New(message)
// }
var PkgNew = errors.New

// New proxy github.com/pkg/errors New
func New(message string, a ...interface{}) error {
	return errors.WithStack(errors.New(message))
}

// Errorf is as the proxy for github.com/pkg/errors.Errorf func.
// func Errorf(format string, args ...interface{}) error {
// 	return errors.Errorf(format, args...)
// }
var Errorf = errors.Errorf

// Wrap is as the proxy for github.com/pkg/errors.Wrap func.
// func Wrap(err error, message string) error {
// 	return errors.Wrap(err, message)
// }
var Wrap = errors.Wrap

// Wrapf is as the proxy for github.com/pkg/errors.Wrapf func.
// func Wrapf(err error, format string, args ...interface{}) error {
// 	return errors.Wrapf(err, format, args...)
// }
var Wrapf = errors.Wrapf

// WithMessage is as the proxy for github.com/pkg/errors.WithMessage func.
// func WithMessage(err error, message string) error {
// 	return errors.WithMessage(err, message)
// }
var WithMessage = errors.WithMessage

// WithMessagef is as the proxy for github.com/pkg/errors.WithMessagef func.
// func WithMessagef(err error, format string, args ...interface{}) error {
// 	return errors.WithMessagef(err, format, args...)
// }
var WithMessagef = errors.WithMessagef

// Cause is as the proxy for github.com/pkg/errors.Cause func.
// func Cause(err error) error {
// 	return errors.Cause(err)
// }
//var Cause = errors.Cause

// WithStack is as the proxy for github.com/pkg/errors.WithStack func.
// func WithStack(err error) error {
// 	return errors.WithStack(err)
// }
var WithStack = errors.WithStack

// Is reports whether any error in err's chain matches target.
// The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.
// An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool such that Is(target) returns true.
var Is = errors.Is

// Cause ...
var Cause = errors.Cause
