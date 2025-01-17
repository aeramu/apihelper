package exception

import "errors"

// NewError creates a new error with the provided options
func NewCustomError(status string, code string, message string, err error) error {
	if err == nil {
		err = errors.New(message)
	}
	return &exception{
		status:  status,
		code:    code,
		message: message,
		error:   err,
	}
}

func newError(status string, message string, err error) error {
	return NewCustomError(status, status, message, err)
}

func InvalidRequest(code, message string, err error) error {
	return NewCustomError(CodeInvalidRequest, code, message, err)
}

func ValidationFailed(code, message string, err error) error {
	return NewCustomError(CodeValidationFailed, code, message, err)
}

func PermissionDenied(code, message string, err error) error {
	return NewCustomError(CodePermissionDenied, code, message, err)
}

func NotFound(code, message string, err error) error {
	return NewCustomError(CodeNotFound, code, message, err)
}

func ThirdParty(code, message string, err error) error {
	return NewCustomError(CodeThirdParty, code, message, err)
}

func AlreadyExists(code, message string, err error) error {
	return NewCustomError(CodeAlreadyExists, code, message, err)
}

func SoftError(code, message string, err error) error {
	return NewCustomError(CodeSoftError, code, message, err)
}

func RaceCondition(code, message string, err error) error {
	return NewCustomError(CodeRaceCondition, code, message, err)
}

func ResourceExhausted(code, message string, err error) error {
	return NewCustomError(CodeResourceExhausted, code, message, err)
}

func Unauthenticated(code, message string, err error) error {
	return NewCustomError(CodeUnauthenticated, code, message, err)
}

func Internal(code, message string, err error) error {
	return NewCustomError(CodeInternal, code, message, err)
}

func Unavailable(code, message string, err error) error {
	return NewCustomError(CodeUnavailable, code, message, err)
}

func DeadlineExceeded(code, message string, err error) error {
	return NewCustomError(CodeDeadlineExceeded, code, message, err)
}
