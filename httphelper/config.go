package httphelper

// Configuration options
type config struct {
	defaultErrorCode    string
	defaultErrorMessage string
	includeDetails      bool
}

const (
	// INTERNAL_SERVER_ERROR is the error code used when an internal server error occurs
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	// INTERNAL_SERVER_MESSAGE provides a descriptive message for internal server errors
	INTERNAL_SERVER_MESSAGE = "An internal server error occurred"
)

// Option represents a configuration option for the httphelper package
type Option func(*config)

// DefaultConfig represents the default configuration
var defaultConfig = config{
	defaultErrorCode:    INTERNAL_SERVER_ERROR,
	defaultErrorMessage: INTERNAL_SERVER_MESSAGE,
	includeDetails:      true,
}

// WithDefaultErrorCode sets the default error code for non-HTTPError errors
func WithDefaultErrorCode(code string) Option {
	return func(c *config) {
		c.defaultErrorCode = code
	}
}

// WithDefaultErrorMessage sets the default error message for non-HTTPError errors
func WithDefaultErrorMessage(msg string) Option {
	return func(c *config) {
		c.defaultErrorMessage = msg
	}
}

// WithIncludeDetails enables or disables including error details in responses
func WithIncludeDetails(include bool) Option {
	return func(c *config) {
		c.includeDetails = include
	}
}

// Configure applies the given options to the package configuration
func Configure(opts ...Option) {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	defaultConfig = cfg
}
