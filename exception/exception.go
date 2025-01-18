package exception

import (
	"net/http"
)

type exception struct {
	s       string
	error   error
	status  string
	code    string
	message string
}

func (e *exception) Error() string {
	return e.s
}

// Code returns the error code
func (e *exception) Code() string {
	return e.code
}

// Message returns the human-readable message
func (e *exception) Message() string {
	return e.message
}

// Unwrap implements the errors.Unwrap interface
func (e *exception) Unwrap() error {
	return e.error
}

// ToHTTPStatus converts an AppError code to an HTTP status code
func (e *exception) HTTPStatus() int {
	switch e.status {
	case CodeInternal:
		return http.StatusInternalServerError // 500
	case CodeInvalidRequest:
		return http.StatusBadRequest // 400
	case CodeValidationFailed:
		return http.StatusUnprocessableEntity // 422
	case CodeNotFound:
		return http.StatusNotFound // 404
	case CodeAlreadyExists, CodeRaceCondition:
		return http.StatusConflict // 409
	case CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodePermissionDenied:
		return http.StatusForbidden // 403
	case CodeResourceExhausted:
		return http.StatusTooManyRequests // 429
	case CodeUnavailable:
		return http.StatusServiceUnavailable // 503
	case CodeDeadlineExceeded:
		return http.StatusGatewayTimeout // 504
	case CodeSoftError:
		return http.StatusOK // 200
	default:
		return http.StatusInternalServerError
	}
}

func (e *exception) GRPCStatus() string {
	switch e.status {
	case CodeInternal:
		return "INTERNAL"
	case CodeInvalidRequest, CodeValidationFailed:
		return "INVALID_ARGUMENT"
	case CodeNotFound:
		return "NOT_FOUND"
	case CodeAlreadyExists, CodeRaceCondition:
		return "ALREADY_EXISTS"
	case CodeUnauthenticated:
		return "UNAUTHENTICATED"
	case CodePermissionDenied:
		return "PERMISSION_DENIED"
	case CodeResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case CodeUnavailable:
		return "UNAVAILABLE"
	case CodeDeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case CodeSoftError:
		return "OK"
	default:
		return "UNKNOWN"
	}
}
