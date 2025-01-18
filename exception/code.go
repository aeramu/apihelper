package exception

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
	ErrorInvalidRequest    = newError(CodeInvalidRequest, "invalid request")
	ErrorValidationFailed  = newError(CodeValidationFailed, "validation failed")
	ErrorPermissionDenied  = newError(CodePermissionDenied, "permission denied")
	ErrorNotFound          = newError(CodeNotFound, "data not found")
	ErrorAlreadyExists     = newError(CodeAlreadyExists, "data already exists")
	ErrorRaceCondition     = newError(CodeRaceCondition, "race condition")
	ErrorResourceExhausted = newError(CodeResourceExhausted, "resource exhausted")
	ErrorUnauthenticated   = newError(CodeUnauthenticated, "unauthenticated")
	ErrorInternal          = newError(CodeInternal, "internal server error")
	ErrorUnavailable       = newError(CodeUnavailable, "service unavailable")
	ErrorDeadlineExceeded  = newError(CodeDeadlineExceeded, "deadline exceeded")
	ErrorSoftError         = newError(CodeSoftError, "soft error")

	// Common errors
)

func newError(status string, message string) error {
	return New(message,
		WithStatus(status),
		WithCode(status),
		WithMessage(message),
	)
}
