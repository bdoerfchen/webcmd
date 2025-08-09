package config

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

const RouteParamPrefix = "WC_"

var allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

type Route struct {
	Method      string
	Route       string
	Command     string
	Args        []string
	StatusCodes []ExitCodeMapping
	Env         map[string]string
}

type ExitCodeMapping struct {
	ExitCode      *int // The base exit code from which to map from
	StatusCode    int  // Status code to map to
	ResponseEmpty bool // Send empty response for this exit code
}

func DefaultRoute() Route {
	defaultCommand := "bash"
	defaultArgs := []string{"-c", "echo This is webcmd. Who are you?"}
	if runtime.GOOS == "windows" {
		defaultCommand = "cmd"
		defaultArgs = []string{"/C", "echo This is webcmd. Who are you?"}
	}

	var zero int = 0
	return Route{
		Method:  http.MethodGet,
		Route:   "/{FILE}",
		Command: defaultCommand,
		Args:    defaultArgs,
		StatusCodes: []ExitCodeMapping{
			{ExitCode: &zero, StatusCode: 200, ResponseEmpty: false},
		},
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
		result = append(result, RouteError{Message: fmt.Sprintf("http method '%s' is now allowed", r.Method), Level: ErrorLevelCritical})
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

	return
}
