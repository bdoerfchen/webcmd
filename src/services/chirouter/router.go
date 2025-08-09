package chirouter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/bdoerfchen/webcmd/src/interfaces"
	"github.com/bdoerfchen/webcmd/src/logging"
	"github.com/bdoerfchen/webcmd/src/model/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type chirouter struct {
	router chi.Router
}

var validMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

func New(ctx context.Context, routes []config.Route, executor interfaces.Executer) *chirouter {
	r := &chirouter{router: chi.NewRouter()}

	// A good base middleware stack
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.StripSlashes)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Recoverer)

	logger := logging.FromContext(ctx)
	logger.Info("begin route registration", slog.Int("count", len(routes)))

	// Register all routes
	for _, route := range routes {
		r.addRoute(route, executor, logger)
	}

	return r
}

func (r *chirouter) addRoute(route config.Route, executor interfaces.Executer, logger *slog.Logger) {
	routeLogger := logger.With(slog.String("route", route.Route))
	routePattern, _ := strings.CutSuffix(route.Route, "/")

	// Check method is valid or skip route
	method := strings.ToUpper(route.Method)
	if !slices.Contains(validMethods, method) {
		logger.Warn("route fail: "+routePattern+" (invalid method for route)",
			slog.String("method", method),
		)
		return
	}

	// Register route to router
	r.router.MethodFunc(method, routePattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startTime := time.Now()

		// Reusable end-of-request logging function
		logFn := func(responseSize int, responseCode int) {
			url := r.URL.String()
			elapsed := time.Since(startTime)
			logger.InfoContext(ctx, fmt.Sprintf("%s %s -> %v", method, url, responseCode),
				slog.Duration("responseTime", elapsed),
				slog.Int("size", responseSize),
				slog.String("userAgent", r.UserAgent()),
			)
		}

		// On handle, start executor for route
		result, exitCode, err := executor.Execute(ctx, route)
		if err != nil {
			// Unexpected error, code 500, no response body
			routeLogger.ErrorContext(ctx,
				"unexpected error while handling route",
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)

			// Log and finish
			logFn(0, http.StatusInternalServerError)
			return
		}

		// Respond with command result and mapped status code from exit code
		responseCode := route.StatusFromExitCode(exitCode)
		w.WriteHeader(responseCode)
		w.Write(result)

		// Log and finish
		logFn(len(result), responseCode)
	})

	logger.Info("route ok: " + routePattern)
}

func (r *chirouter) Handler() http.Handler {
	return r.router
}
