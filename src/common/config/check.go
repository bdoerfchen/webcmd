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

	// Check caching
	if r.Caching && r.Method != http.MethodGet {
		result = append(result, RouteError{Message: "caching only works on GET requests", Level: ErrorLevelWarning})
	}

	// Check exec
	if r.Exec.Proc == nil && r.Exec.Shell == nil {
		result = append(result, RouteError{Message: "exec requires 'proc' or 'shell' config", Level: ErrorLevelCritical})
	} else if r.Exec.Proc != nil && r.Exec.Shell != nil {
		result = append(result, RouteError{Message: "'shell' config will be ignored when providing 'proc' config", Level: ErrorLevelWarning})
	}

	// Check exec.proc
	if r.Exec.Proc != nil {
		if r.Exec.Proc.Path == "" {
			result = append(result, RouteError{Message: "executable path must not be empty", Level: ErrorLevelCritical})
		} else if _, err := exec.LookPath(r.Exec.Proc.Path); err != nil {
			result = append(result, RouteError{Message: fmt.Sprintf("executable '%s' can not be found as file or on PATH", r.Exec.Proc.Path), Level: ErrorLevelWarning})
		}
	} else if r.Exec.Shell != nil {
		if r.Exec.Shell.Command == "" {
			result = append(result, RouteError{Message: "shell command must not be empty", Level: ErrorLevelCritical})
		}
	}

	// Check exit codes
	for _, codeMapping := range r.StatusCodes {
		// Check for valid response stream names
		if !codeMapping.ResponseStream.IsValid() {
			result = append(result, RouteError{Message: fmt.Sprintf("exit code %v with invalid response stream '%s'", *codeMapping.ExitCode, codeMapping.ResponseStream), Level: ErrorLevelCritical})
		}
	}

	// Check default status code
	if !slices.ContainsFunc(r.StatusCodes, func(i ExitCodeMapping) bool { return i.ExitCode == nil }) {
		result = append(result, RouteError{Message: "no default status code for non-zero exit codes defined: uses 500 now", Level: ErrorLevelInfo})
	}

	// Check query params
	for _, param := range r.QueryParams {
		if invalidQueryParam.MatchString(param) {
			result = append(result, RouteError{Message: fmt.Sprintf("query param '%s' is not a valid name", param), Level: ErrorLevelCritical})
		}
	}

	// OS specific remarks
	if runtime.GOOS == "windows" {
		if r.Exec.Shell != nil {
			result = append(result, RouteError{Message: "the 'shell' exec mode is not supported on windows", Level: ErrorLevelCritical})
		}
	}

	return
}

// Must not start with a number, or contain something else than ASCII characters, numbers or underscores
var invalidQueryParam = regexp.MustCompile(`^[\d]|[^\w\d_]+`)
