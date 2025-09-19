package chirouter

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"

	"github.com/bdoerfchen/webcmd/src/common/config"
	"github.com/bdoerfchen/webcmd/src/common/process"
	"github.com/go-chi/chi/v5"
)

const DefaultKey = -1

type OptimizedRoute struct {
	config.Route
	StatusCodeMap map[int]OptimizedMapping
	ParamNames    []string
}

type OptimizedMapping struct {
	config.ExitCodeMapping
}

func OptimizeRoute(route config.Route) (result OptimizedRoute) {
	result.Route = route
	result.StatusCodeMap = make(map[int]OptimizedMapping)

	// Convert all mappings and add them to map
	for _, codeMap := range route.StatusCodes {
		key := DefaultKey
		if codeMap.ExitCode != nil {
			key = *codeMap.ExitCode
		}

		if codeMap.StatusCode < http.StatusOK {
			// Server should not answer below 200 when request is finished
			codeMap.StatusCode = http.StatusOK
		}
		if codeMap.StatusCode >= 1000 {
			// Arbitrary limit
			codeMap.StatusCode = 999
		}

		// Use default output stream if left empty
		if codeMap.ResponseStream == "" {
			codeMap.ResponseStream = route.ResponseStream
		}

		// Add entry
		result.StatusCodeMap[key] = OptimizedMapping{codeMap}
	}

	// Add defaults
	if _, ok := result.StatusCodeMap[DefaultKey]; !ok {
		result.StatusCodeMap[DefaultKey] = OptimizedMapping{
			config.ExitCodeMapping{
				StatusCode:     http.StatusInternalServerError,
				ResponseStream: config.Both,
			},
		}
	}

	// Parse paramters
	result.ParamNames = paramNamesInRoute(route.Route)

	return
}

func (o *OptimizedRoute) ExitCodeResponse(code int) OptimizedMapping {
	if response, ok := o.StatusCodeMap[code]; ok {
		return response
	}

	return o.StatusCodeMap[DefaultKey]
}

func (o *OptimizedMapping) ResponseBufferFor(proc *process.Process) *bytes.Buffer {
	if o == nil {
		return nil
	}

	switch strings.ToLower(string(o.ExitCodeMapping.ResponseStream)) {
	case string(config.StdOut):
		return &proc.StdOut
	case string(config.StdErr):
		return &proc.StdErr
	case string(config.None):
		return nil
	default:
		return &proc.StdOutErr
	}
}

var paramMatcher = regexp.MustCompile(`{([^\/\\:]+)(?:\:[^}\/]+)*}`)

func paramNamesInRoute(routePattern string) (result []string) {
	groups := paramMatcher.FindAllStringSubmatch(routePattern, -1)
	if groups == nil {
		return
	}

	for _, g := range groups {
		result = append(result, g[1])
	}

	return
}

// Read all query and path parameters of a http.Request. Path params are dominant
func (o *OptimizedRoute) RequestParameters(r *http.Request) map[string]string {
	result := make(map[string]string)

	// Add query parameters
	for _, paramName := range o.QueryParams {
		result[paramName] = r.URL.Query().Get(paramName)
	}

	// Add route parameters, potentially overwriting query params on double initialization
	for _, paramName := range o.ParamNames {
		result[paramName] = chi.URLParam(r, paramName)
	}

	return result
}
