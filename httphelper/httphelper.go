package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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


// HTTPError defines the interface for custom error types that provide
// detailed HTTP error responses. Implementing this interface allows
// errors to control their HTTP representation.
type HTTPError interface {
	// Error returns the complete error message
	Error() string
	// HTTPStatus returns the appropriate HTTP status code for the error
	HTTPStatus() int
	// Message returns a human-readable error description
	Message() string
	// Code returns a machine-readable error identifier
	Code() string
}

func AsHTTPError(err error) (HTTPError, bool) {
	if err == nil {
		return nil, false
	}

	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
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
	var httpStatus int
	if httpErr, ok := AsHTTPError(err); ok {
		errInfo = ErrorInfo{
			Code:    httpErr.Code(),
			Message: httpErr.Message(),
		}
		if defaultConfig.includeDetails {
			errInfo.Detail = httpErr.Error()
		}
		httpStatus = httpErr.HTTPStatus()
	} else {
		errInfo = ErrorInfo{
			Code:    defaultConfig.defaultErrorCode,
			Message: defaultConfig.defaultErrorMessage,
		}
		if defaultConfig.includeDetails {
			errInfo.Detail = err.Error()
		}
		httpStatus = http.StatusInternalServerError
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(Response{
		Status:    httpStatus,
		Success:   false,
		ErrorInfo: &errInfo,
	})
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
	if err := r.Err(); err != nil {
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
