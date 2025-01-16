package exception

func (e *codeError) GRPCStatus() string {
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
