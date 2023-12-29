package webbase

import (
	"fmt"
	"net/http"
)

type serveConfiguration struct {
	webListenAddress       string
	serviceListenAddress   string
	sentryDebug            bool
	enableServiceListener  bool
	healthCheckHandlerFunc http.HandlerFunc
}

type serveOption func(*serveConfiguration) error

// WithWebListenAddress sets the address for the web listener
func WithWebListenAddress(address string) serveOption {
	return func(c *serveConfiguration) error {
		if address == "" {
			return fmt.Errorf("webListenAddress must not be empty")
		}
		c.webListenAddress = address
		return nil
	}
}

// WithServiceListenAddress sets the address for the service (healthcheck/metrics) listener
func WithServiceListenAddress(address string) serveOption {
	return func(c *serveConfiguration) error {
		if address == "" {
			return fmt.Errorf("serviceListenAddress must not be empty")
		}
		c.serviceListenAddress = address
		return nil
	}
}

// WithSentryDebug sets the debug flag for sentry
func WithSentryDebug(debug bool) serveOption {
	return func(c *serveConfiguration) error {
		c.sentryDebug = debug
		return nil
	}
}

// WithoutServiceEndpoint disables the metrics/healthcheck endpoint
func WithoutServiceEndpoint() serveOption {
	return func(c *serveConfiguration) error {
		c.enableServiceListener = false
		return nil
	}
}

// WithHealthCheckHandlerFunc sets the handler function for the healthcheck endpoint
func WithHealthCheckHandlerFunc(handlerFunc http.HandlerFunc) serveOption {
	return func(c *serveConfiguration) error {
		if handlerFunc == nil {
			return fmt.Errorf("healthCheckHandler must not be nil")
		}
		c.healthCheckHandlerFunc = handlerFunc
		return nil
	}
}
