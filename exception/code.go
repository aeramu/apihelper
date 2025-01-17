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