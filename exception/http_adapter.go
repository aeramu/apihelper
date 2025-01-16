package exception

import "net/http"

// ToHTTPStatus converts an AppError code to an HTTP status code
func (e *codeError) HTTPCode() int {
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
