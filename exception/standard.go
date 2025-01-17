package exception

import "errors"

// ErrorCode represents an error with additional context through code and message.
type ErrorCode interface {
	Error() string
	Code() string
}

// AsErrorCode attempts to convert an error to a ErrorCode interface.
// If the error already implements ErrorCode, it returns the error as is.
// If not, it returns nil.
func AsErrorCode(err error) (ErrorCode, bool) {
	if err == nil {
		return nil, false
	}

	var errorCode ErrorCode
	if errors.As(err, &errorCode) {
		return errorCode, true
	}
	return nil, false
}