package chirouter

import (
	"net/http"
	"regexp"

	"github.com/bdoerfchen/webcmd/src/model/config"
	"github.com/go-chi/chi/v5"
)

const DefaultKey = -1

type OptimizedRoute struct {
	config.Route
	StatusCodeMap map[int]config.ExitCodeMapping
	ParamNames    []string
}

func OptimizeRoute(route config.Route) (result OptimizedRoute) {
	result.Route = route
	result.StatusCodeMap = make(map[int]config.ExitCodeMapping)

	// Convert all mapings and add them to map
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

		// Add entry
		result.StatusCodeMap[key] = codeMap
	}

	// Add defaults
	if _, ok := result.StatusCodeMap[DefaultKey]; !ok {
		result.StatusCodeMap[DefaultKey] = config.ExitCodeMapping{
			StatusCode:    http.StatusInternalServerError,
			ResponseEmpty: true,
		}
	}

	// Parse paramters
	result.ParamNames = paramNamesInRoute(route.Route)

	return
}

func (o *OptimizedRoute) ExitCodeResponse(code int) config.ExitCodeMapping {
	if response, ok := o.StatusCodeMap[code]; ok {
		return response
	}

	return o.StatusCodeMap[DefaultKey]
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
