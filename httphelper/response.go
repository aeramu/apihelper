package httphelper

const (
	
	// UNKNOWN_ERROR is the error code used when the error type cannot be determined
	UNKNOWN_ERROR = "UNKNOWN_ERROR"
	// UNKNOWN_DETAIL provides a descriptive message for unknown errors
	UNKNOWN_DETAIL = "Request failed without error details. This may be due to malformed JSON, invalid JSON format, or empty response body"
)

// Response represents a standardized API response structure.
// It provides a consistent format for both successful and error responses,
// making it easier to handle API responses in a uniform way.
type Response struct {
	// Status represents the HTTP status code of the response (e.g., 200, 404, 500)
	Status int `json:"status"`
	// Success indicates whether the request was processed successfully
	Success bool `json:"success"`
	// Data contains the response payload for successful requests
	// For error responses, this field will be null
	Data any `json:"data"`
	// ErrorInfo contains error details when Success is false
	// This field is omitted for successful responses
	ErrorInfo *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo provides structured error information for API responses.
// It is designed to give both human-readable and machine-processable
// error details.
type ErrorInfo struct {
	// Code is a machine-readable identifier for the error type
	// Examples: "INVALID_INPUT", "RESOURCE_NOT_FOUND"
	Code string `json:"code"`
	// Message is a human-readable description of the error
	Message string `json:"message"`
	// Detail is a technical description of the error (optional)
	// This provides more specific information about what went wrong
	Detail string `json:"detail,omitempty"`
	// Details contains additional error context (optional)
	// This can be structured data providing more information about the error
	Details any `json:"details,omitempty"`
}

func (r *Response) IsSuccess() bool {
	return r.Success
}

func (r *Response) IsError() bool {
	return !r.Success
}

func (r *Response) HTTPStatus() int {
	return r.Status
}

func (r *Response) Error() string {
	err := r.getError()
	if err == nil {
		return ""
	}
	return err.Detail
}

func (r *Response) Code() string {
	err := r.getError()
	if err == nil {
		return ""
	}
	return err.Code
}

func (r *Response) Message() string {
	err := r.getError()
	if err == nil {
		return ""
	}
	return err.Message
}

func (r *Response) Err() error {
	if r.IsSuccess() {
		return nil
	}
	err := r.getError()
	if err == nil {
		return nil
	}
	r.ErrorInfo = err
	return r
}

func (r *Response) getError() *ErrorInfo {
	if r.IsSuccess() {
		return nil
	}
	if r.ErrorInfo == nil || (r.ErrorInfo.Code == "" && r.ErrorInfo.Detail == "") {
		return &ErrorInfo{
			Code:   UNKNOWN_ERROR,
			Detail: UNKNOWN_DETAIL,
		}
	}
	return r.ErrorInfo
}
