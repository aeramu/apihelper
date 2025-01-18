package exception

import (
	"fmt"
)

// ErrorOption is a function that configures an error
type ErrorOption func(*exception)

// defaultOptions stores the global default options
var defaultOptions []ErrorOption = []ErrorOption{
	WithStatus(CodeInternal),
	WithCode(CodeInternal),
}

// SetDefaultOptions sets the global default options that will be applied to all new exceptions
func SetDefaultOptions(opts ...ErrorOption) {
	defaultOptions = opts
}

func WithError(err error) ErrorOption {
	return func(e *exception) {
		if err == nil {
			return
		}
		e.error = err
	}
}

func WithStatus(status string) ErrorOption {
	return func(e *exception) {
		e.status = status
	}
}

func WithCode(code string) ErrorOption {
	return func(e *exception) {
		e.code = code
	}
}

func WithMessage(message string) ErrorOption {
	return func(e *exception) {
		e.message = message
	}
}

func WithArgs(args ...any) ErrorOption {
	return func(e *exception) {
		e.s = fmt.Sprintf(e.s, args...)
	}
}

// New creates a new Exception with required code and message, plus optional configurations
func New(text string, opts ...ErrorOption) error {
	e := &exception{
		s: text,
	}

	// Apply default options first
	for _, opt := range defaultOptions {
		opt(e)
	}

	// Apply instance-specific options, which can override defaults
	for _, opt := range opts {
		opt(e)
	}

	if e.error != nil {
		e.s = fmt.Sprintf("%s: %s", e.s, e.error.Error())
	}

	return e
}
