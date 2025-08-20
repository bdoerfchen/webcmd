package config

import (
	"net/http"
)

const RouteParamPrefix = "WC_"

var allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

type Route struct {
	Method      string            // HTTP method
	Route       string            // Route pattern, includes the url path parameters
	Headers     map[string]string // Default response headers for this route
	QueryParams []string          // List of allowed url query parameters
	StatusCodes []ExitCodeMapping // List of exit-code to status-code mappings
	Env         map[string]string // Environment variable map
	AllowBody   bool              // Enable reading the request body and writing it into stdin of the exec environment
	Exec        RouteExec
}

type ExitCodeMapping struct {
	ExitCode      *int              // The base exit code from which to map from
	StatusCode    int               // Status code to map to
	ResponseEmpty bool              // Send empty response for this exit code
	Headers       map[string]string // Special response headers for this exit code
}

// Return a default route configuration that can be used as the base for further configuration.
func DefaultRoute() Route {
	var zero int = 0
	return Route{
		Method: http.MethodGet,
		Route:  "/*",
		Exec:   RouteExec{},
		StatusCodes: []ExitCodeMapping{
			{ExitCode: &zero, StatusCode: 200},
		},
		Env:         make(map[string]string),
		QueryParams: make([]string, 0),
	}
}

// Prints "METHOD PATH" (example: GET /)
func (r *Route) String() string {
	return r.Method + " " + r.Route
}
