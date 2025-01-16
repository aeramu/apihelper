package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	// UNKNOWN_ERROR is the error code used when the error type cannot be determined
	UNKNOWN_ERROR = "UNKNOWN_ERROR"
	// UNKNOWN_MESSAGE provides a descriptive message for unknown errors
	UNKNOWN_MESSAGE = "Request failed without error details. This may be due to malformed JSON, invalid JSON format, or empty response body"
)

// Package httphelper provides utilities for standardized HTTP response handling.
// It offers a consistent way to structure API responses, handle errors, and
// safely extract typed data from responses.
//
// Example usage:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    data := map[string]string{"message": "Hello World"}
//	    httphelper.OK(w, data)
//	}
//
//	func errorHandler(w http.ResponseWriter, r *http.Request) {
//	    err := errors.New("something went wrong")
//	    httphelper.Error(w, err)
//	}

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

// GetErrorInfo returns error information from the response.
// Returns:
//   - nil if the response was successful (Success = true)
//   - ErrorInfo containing error details if Success = false
//   - A generic ErrorInfo if Success = false but no error details were provided
func (r Response) GetErrorInfo() *ErrorInfo {
	if r.Success {
		return nil
	}
	if r.ErrorInfo == nil || (r.ErrorInfo.Code == "" && r.ErrorInfo.Message == "") {
		return &ErrorInfo{
			Code:    UNKNOWN_ERROR,
			Message: UNKNOWN_MESSAGE,
		}
	}
	return r.ErrorInfo
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
	// Details contains additional error context (optional)
	// This can be structured data providing more information about the error
	Details any `json:"details,omitempty"`
}

// Error implements the error interface for ErrorInfo.
// It formats the error message in a consistent way:
// - If ErrorInfo is nil, returns "error is nil"
// - If Details are present, returns "message (code): details"
// - If Code is empty, returns just the message
// - Otherwise returns "message (code)"
func (e *ErrorInfo) Error() string {
	if e == nil {
		return ""
	}
	if e.Details != nil {
		return fmt.Sprintf("%s (%s): %v", e.Message, e.Code, e.Details)
	}
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// ReadData safely extracts and unmarshals the response Data field into the specified type T.
// It handles various data formats and provides type-safe data extraction.
//
// Example usage:
//
//	var users []User
//	data, err := ReadData[[]User](response)
//	if err != nil {
//	    return fmt.Errorf("failed to read users: %w", err)
//	}
//
// Parameters:
//   - r: The Response object containing the data to extract
//
// Returns:
//   - The unmarshaled data of type T
//   - An error if the response contains an error or if unmarshaling fails
func ReadData[T any](r Response) (T, error) {
	var data T
	// First check if response is successful
	if err := r.GetErrorInfo(); err != nil {
		return data, err
	}

	// Handle nil data
	if r.Data == nil {
		return data, fmt.Errorf("response data is nil")
	}

	// Convert data to JSON bytes for consistent unmarshaling
	var jsonBytes []byte
	switch v := r.Data.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		var err error
		jsonBytes, err = json.Marshal(r.Data)
		if err != nil {
			return data, fmt.Errorf("failed to marshal response data: %w", err)
		}
	}

	// Unmarshal JSON bytes into target type
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return data, nil
}

// HTTPError defines the interface for custom error types that provide
// detailed HTTP error responses. Implementing this interface allows
// errors to control their HTTP representation.
type HTTPError interface {
	// Error returns the complete error message
	Error() string
	// HTTPCode returns the appropriate HTTP status code for the error
	HTTPCode() int
	// Message returns a human-readable error description
	Message() string
	// Code returns a machine-readable error identifier
	Code() string
}

// OK writes a successful JSON response with the provided data.
// It automatically sets the appropriate headers and status code.
//
// Parameters:
//   - w: The HTTP response writer
//   - data: The data to include in the response
func OK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Success: true,
		Data:    data,
	})
}

// Error writes an error response in JSON format.
// It handles both standard errors and custom errors implementing the HTTPError interface.
//
// Parameters:
//   - w: The HTTP response writer
//   - err: The error to include in the response
func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var errInfo ErrorInfo
	var httpCode int
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		errInfo = ErrorInfo{
			Code:    httpErr.Code(),
			Message: httpErr.Message(),
		}
		if defaultConfig.includeDetails {
			errInfo.Details = httpErr.Error()
		}
		httpCode = httpErr.HTTPCode()
	} else {
		errInfo = ErrorInfo{
			Code:    defaultConfig.defaultErrorCode,
			Message: defaultConfig.defaultErrorMessage,
		}
		if defaultConfig.includeDetails {
			errInfo.Details = err.Error()
		}
		httpCode = http.StatusInternalServerError
	}

	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(Response{
		Status:    httpCode,
		Success:   false,
		ErrorInfo: &errInfo,
	})
}
