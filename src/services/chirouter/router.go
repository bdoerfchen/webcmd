package chirouter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bdoerfchen/webcmd/src/common/config"
	"github.com/bdoerfchen/webcmd/src/common/execution"
	"github.com/bdoerfchen/webcmd/src/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type chirouter struct {
	router             chi.Router
	sanitizer          *valueSanitizer
	executerCollection *execution.ExecuterCollection
}

func New(executerCollection *execution.ExecuterCollection) *chirouter {
	return &chirouter{
		router:             chi.NewRouter(),
		sanitizer:          newSanitizer(),
		executerCollection: executerCollection,
	}
}

func (r *chirouter) Register(ctx context.Context, routes []config.Route) error {
	// Basic middleware registration
	r.router.Use(middleware.StripSlashes)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Recoverer)

	logger := logging.FromContext(ctx)
	logger.Debug("begin route registration", slog.Int("count", len(routes)))

	// Register all routes
	for _, route := range routes {
		executer, err := r.executerCollection.For(&route)
		if err != nil {
			logger.Error("no executer available for route " + route.String())
			continue
		}
		r.addRoute(route, executer, logger)
	}

	logger.Debug("route registration done")

	return nil
}

// Actual route registration with the handler function definition
func (r *chirouter) addRoute(route config.Route, executor execution.Executer, logger *slog.Logger) {
	routeLogger := logger.With(slog.String("route", route.Route))
	routePattern, _ := strings.CutSuffix(route.Route, "/")

	// Reusable end-of-request logging function for this route
	logFn := func(startTime time.Time, req *http.Request, responseSize int, responseCode int) {
		url := req.URL.String()
		elapsed := time.Since(startTime)
		logger.Info(fmt.Sprintf("%s %s -> %v", route.Method, url, responseCode),
			slog.Duration("responseTime", elapsed),
			slog.Int("size", responseSize),
			slog.String("userAgent", req.UserAgent()),
		)
	}

	// Register route to router
	optimizedRoute := OptimizeRoute(route)
	r.router.MethodFunc(route.Method, routePattern, func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		startTime := time.Now()

		// Execution config
		execConfig := execution.ConfigFromRoute(&route)
		if route.AllowBody {
			execConfig.Stdin = req.Body
		}

		// Load URL parameters as env variables
		params := optimizedRoute.RequestParameters(req)
		for key, value := range params {
			// Skip unset values to keep defaults and reduce sanitization efforts
			if value == "" {
				continue
			}
			// convert key to env variable format
			key = config.RouteParamPrefix + strings.ToUpper(key)
			// Sanitize input and add to route env map
			execConfig.Env[key] = r.sanitizer.Sanitize(value)
		}

		// On handle, start executor for route
		result, exitCode, err := executor.Execute(ctx, execConfig)
		if err != nil {
			// Unexpected error, code 500, no response body
			routeLogger.ErrorContext(ctx,
				"unexpected error while handling route",
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)

			// Log and finish
			logFn(startTime, req, 0, http.StatusInternalServerError)
			return
		}

		// Load response config for exit code
		exitResponse := optimizedRoute.ExitCodeResponse(exitCode)

		// Set headers (default and exit code related)
		for header, value := range route.Headers {
			w.Header().Add(header, value)
		}
		for header, value := range exitResponse.Headers {
			w.Header().Add(header, value)
		}

		// Respond with command result and mapped status code from exit code
		w.WriteHeader(exitResponse.StatusCode)
		writtenLen := 0
		if !exitResponse.ResponseEmpty {
			w.Write(result.StdOutErr.Bytes())
			writtenLen = len(result.StdOutErr.Bytes())
		}

		// Log and finish
		logFn(startTime, req, writtenLen, exitResponse.StatusCode)
	})

	logger.Debug(fmt.Sprintf("- %s %s", route.Method, routePattern))
}

func (r *chirouter) Handler() http.Handler {
	return r.router
}
