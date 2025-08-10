package config

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

const RouteParamPrefix = "WC_"

var allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

type Route struct {
	Method      string            // HTTP method
	Route       string            // Route pattern, includes the url path parameters
	Command     string            // Command path
	Args        []string          // List of arguments for the command
	QueryParams []string          // List of allowed url query parameters
	StatusCodes []ExitCodeMapping // List of exit-code to status-code mappings
	Env         map[string]string // Environment variable map
	Headers     map[string]string // Default response headers for this route
	AllowBody   bool              // Enable reading the request body and writing it into stdin of the exec environment
}

type ExitCodeMapping struct {
	ExitCode      *int              // The base exit code from which to map from
	StatusCode    int               // Status code to map to
	ResponseEmpty bool              // Send empty response for this exit code
	Headers       map[string]string // Special response headers for this exit code
}

// Return a default route configuration that can be used as the base for further configuration.
func DefaultRoute() Route {
	defaultCommand := "bash"
	const defaultMessage = "echo Your webcmd works! Visit https://github.com/bdoerfchen/webcmd to learn more about how to use it."
	defaultArgs := []string{"-c", defaultMessage}
	if runtime.GOOS == "windows" {
		defaultCommand = "cmd"
		defaultArgs = []string{"/C", defaultMessage}
	}

	var zero int = 0
	return Route{
		Method:  http.MethodGet,
		Route:   "/*",
		Command: defaultCommand,
		Args:    defaultArgs,
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

// Perform check on all fields and return a collection of remarks
func (r *Route) Check() (result RouteErrorCollection) {
	// Check method
	r.Method = strings.ToUpper(r.Method)
	if !slices.Contains(allowedMethods, r.Method) {
		result = append(result, RouteError{Message: fmt.Sprintf("http method '%s' is not allowed", r.Method), Level: ErrorLevelCritical})
	}

	// Add info when body is not allowed for POST or PUT route
	if (r.Method == http.MethodPost || r.Method == http.MethodPut) && !r.AllowBody {
		result = append(result, RouteError{Message: "body will be ignored", Level: ErrorLevelInfo})
	}

	// Check command
	if r.Command == "" {
		result = append(result, RouteError{Message: "command must not be empty", Level: ErrorLevelCritical})
	} else if _, err := exec.LookPath(r.Command); err != nil {
		result = append(result, RouteError{Message: fmt.Sprintf("command '%s' can not be found on PATH", r.Command), Level: ErrorLevelWarning})
	}

	// Check status codes
	if !slices.ContainsFunc(r.StatusCodes, func(i ExitCodeMapping) bool { return i.ExitCode == nil }) {
		result = append(result, RouteError{Message: "no default status code defined", Level: ErrorLevelInfo})
	}

	// Check query params
	for _, param := range r.QueryParams {
		if invalidQueryParam.MatchString(param) {
			result = append(result, RouteError{Message: fmt.Sprintf("query param '%s' is not a valid name", param), Level: ErrorLevelCritical})
		}
	}

	return
}

// Must not start with a number, or contain something else than ASCII characters, numbers or underscores
var invalidQueryParam = regexp.MustCompile(`^[\d]|[^\w\d_]+`)
