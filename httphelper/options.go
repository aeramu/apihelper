package httphelper

// Configuration options
type config struct {
	defaultErrorCode    string
	defaultErrorMessage string
	includeDetails      bool
}

// Option represents a configuration option for the httphelper package
type Option func(*config)

// DefaultConfig represents the default configuration
var defaultConfig = config{
	defaultErrorCode:    "INTERNAL_SERVER_ERROR",
	defaultErrorMessage: "An internal server error occurred",
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

// WithErrorDetails enables or disables including error details in responses
func WithErrorDetails(include bool) Option {
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