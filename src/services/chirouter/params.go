package chirouter

import (
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/bdoerfchen/webcmd/src/common/config"
	"github.com/go-chi/chi/v5"
)

type parameterCollection struct {
	parameters []config.RouteParameter
	sanitizer  *valueSanitizer
}

func NewParameterCollection(route config.Route) *parameterCollection {
	// Copy parameters list and add missing route parameters
	parameters := slices.Clone(route.Parameters)
	for _, param := range paramNamesInRoute(route.Route) {
		if !slices.ContainsFunc(parameters, func(p config.RouteParameter) bool { return p.Name == param }) {
			parameters = append(parameters, config.RouteParameter{
				Name:   param,
				Source: config.ParamSourceRoute,
			})
		}
	}

	return &parameterCollection{
		parameters: parameters,
		sanitizer:  newSanitizer(),
	}
}

// Read all parameters from an http.Request in the order defined in the route. Later params overwrite earlier ones
func (c *parameterCollection) For(r *http.Request) map[string]string {
	result := make(map[string]string)

	// Iterate over all parameters
	for _, param := range c.parameters {
		var value string
		switch param.Source {
		case config.ParamSourceHeader:
			value = r.Header.Get(param.Name)
		case config.ParamSourceQuery:
			value = r.URL.Query().Get(param.Name)
		case config.ParamSourceRoute:
			value = chi.URLParam(r, param.Name)
		}

		if value != "" {
			// Clean user input
			value = c.sanitizer.Sanitize(value)
		} else {
			// Use default if empty
			value = param.Default
		}

		// Define name of env variable
		envName := param.As
		if envName == "" {
			envName = config.RouteParamPrefix + strings.ToUpper(param.Name)
		}

		result[envName] = value
	}

	return result
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
