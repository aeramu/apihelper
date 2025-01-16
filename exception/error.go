package exception

import (
	"errors"
)

// Common error codes that are protocol-agnostic
const (
	// System/Infrastructure errors
	CodeInternal         = "INTERNAL"
	CodeUnavailable      = "UNAVAILABLE"
	CodeDeadlineExceeded = "DEADLINE_EXCEEDED"
	CodeThirdParty       = "THIRD_PARTY"

	// Input/Validation errors
	CodeInvalidRequest   = "INVALID_REQUEST"
	CodeValidationFailed = "VALIDATION_FAILED"

	// Authentication/Authorization errors
	CodeUnauthenticated  = "UNAUTHENTICATED"
	CodePermissionDenied = "PERMISSION_DENIED"

	// Resource errors
	CodeNotFound          = "NOT_FOUND"
	CodeAlreadyExists     = "ALREADY_EXISTS"
	CodeRaceCondition     = "RACE_CONDITION"
	CodeResourceExhausted = "RESOURCE_EXHAUSTED"

	// Soft errors
	CodeSoftError = "SOFT_ERROR"
)

var (
	// Base errors
	ErrorInvalidRequest    = newError(CodeInvalidRequest, "Invalid request", nil)
	ErrorValidationFailed  = newError(CodeValidationFailed, "Validation failed", nil)
	ErrorPermissionDenied  = newError(CodePermissionDenied, "Permission denied", nil)
	ErrorNotFound          = newError(CodeNotFound, "Resource not found", nil)
	ErrorAlreadyExists     = newError(CodeAlreadyExists, "Resource already exists", nil)
	ErrorRaceCondition     = newError(CodeRaceCondition, "Race condition", nil)
	ErrorResourceExhausted = newError(CodeResourceExhausted, "Resource exhausted", nil)
	ErrorUnauthenticated   = newError(CodeUnauthenticated, "Unauthenticated", nil)
	ErrorInternal          = newError(CodeInternal, "Internal server error", nil)
	ErrorUnavailable       = newError(CodeUnavailable, "Service unavailable", nil)
	ErrorDeadlineExceeded  = newError(CodeDeadlineExceeded, "Deadline exceeded", nil)
	ErrorSoftError         = newError(CodeSoftError, "Soft error", nil)

	// Common errors
)

type codeError struct {
	error
	status  string
	code    string
	message string
}

func (e *codeError) Error() string {
	return e.error.Error()
}

// Code returns the error code
func (e *codeError) Code() string {
	return e.code
}

// Message returns the human-readable message
func (e *codeError) Message() string {
	return e.message
}

// Unwrap implements the errors.Unwrap interface
func (e *codeError) Unwrap() error {
	return e.error
}

// NewError creates a new error with the provided options
func NewCustomError(status string, code string, message string, err error) error {
	if err == nil {
		err = errors.New(status)
	}
	return &codeError{
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
