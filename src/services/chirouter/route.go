package chirouter

import (
	"net/http"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

const DefaultKey = -1

type OptimizedRoute struct {
	config.Route
	StatusCodeMap map[int]ExitCodeResponse
}

type ExitCodeResponse struct {
	StatusCode    int
	ResponseEmpty bool
}

func OptimizeRoute(route config.Route) (result OptimizedRoute) {
	result.Route = route
	result.StatusCodeMap = make(map[int]ExitCodeResponse)

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
			// Arbitrary limit - not that I really care
			codeMap.StatusCode = 999
		}

		// Add entry
		result.StatusCodeMap[key] = ExitCodeResponse{
			StatusCode:    codeMap.StatusCode,
			ResponseEmpty: codeMap.ResponseEmpty,
		}
	}

	// Add defaults
	if _, ok := result.StatusCodeMap[DefaultKey]; !ok {
		result.StatusCodeMap[DefaultKey] = ExitCodeResponse{
			StatusCode:    http.StatusInternalServerError,
			ResponseEmpty: true,
		}
	}

	return
}

func (o *OptimizedRoute) ExitCodeResponse(code int) ExitCodeResponse {
	if response, ok := o.StatusCodeMap[code]; ok {
		return response
	}

	return o.StatusCodeMap[DefaultKey]
}
