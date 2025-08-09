package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bdoerfchen/webcmd/src/logging"
	"github.com/bdoerfchen/webcmd/src/model/config"
	"github.com/bdoerfchen/webcmd/src/services/chirouter"
	"github.com/bdoerfchen/webcmd/src/services/executer"
	"github.com/bdoerfchen/webcmd/src/services/server"
)

func main() {
	appConfig := config.DefaultAppConfig()
	appConfig.Routes = []config.Route{
		{
			Route:   "/test/{id}",
			Command: "bash",
			Args:    []string{"-c", "echo $test"},
			Method:  http.MethodGet,
			Env: map[string]string{
				"test": "Hello, this is a test!",
			},
		},
	}

	// Configure logger
	logger := logging.New(slog.LevelDebug)

	// Setup router with cmd executer
	setupCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	setupCtx = logging.AddToContext(setupCtx, logger)
	exec := executer.New()
	router := chirouter.New(setupCtx, appConfig.Routes, exec)

	// Run server
	runCtx, _ := signal.NotifyContext(context.Background(), os.Kill)
	runCtx = logging.AddToContext(runCtx, logger)
	server := server.New(appConfig.ServerConfig)
	server.Run(runCtx, router.Handler())
}
