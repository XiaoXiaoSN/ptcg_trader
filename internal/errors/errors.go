package errors

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

type _error struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func (e *_error) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[%s] %s", e.ErrorCode, e.Message)
	return buf.String()
}

// Is compare are two errors same
func (e *_error) Is(target error) bool {
	causeErr := errors.Cause(target)
	tErr, ok := causeErr.(*_error)
	if !ok {
		return false
	}
	return e.StatusCode == tErr.StatusCode
}

// HTTPError define error http need
type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetHTTPError convert error to http need
func GetHTTPError(err error) HTTPError {
	_err, ok := err.(*_error)
	if !ok || _err == nil {
		_err = ErrInternalError
	}

	return HTTPError{
		Status:  _err.StatusCode,
		Message: _err.Message,
		Code:    _err.ErrorCode,
	}
}
