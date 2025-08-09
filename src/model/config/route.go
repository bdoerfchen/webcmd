package config

import (
	"net/http"
	"strconv"
)

const defaultExitCodeKey = "default"

type Route struct {
	Route       string
	Command     string
	Args        []string
	Method      string
	StatusCodes map[string]int
	Env         map[string]string
}

func (r Route) StatusFromExitCode(exitCode int) int {
	// Get status code directly from exit code
	statusCode, ok := r.StatusCodes[strconv.Itoa(exitCode)]
	if ok {
		return statusCode
	}

	// Get user-defined default status code if not found previously
	statusCode, ok = r.StatusCodes[defaultExitCodeKey]
	if ok {
		return statusCode
	}

	// Return hard-coded default
	if exitCode == 0 {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}
