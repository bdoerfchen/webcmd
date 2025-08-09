package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bdoerfchen/webcmd/src/logging"
)

type server struct {
	config Config
}

func New(config Config) *server {
	return &server{
		config: config,
	}
}

func (s *server) Run(ctx context.Context, handler http.Handler) error {
	// Setup
	logger := logging.FromContext(ctx)
	host := fmt.Sprintf("%s:%v", s.config.Host, s.config.Port)
	logger.Info(fmt.Sprintf("listening on %s", host))

	// Listen
	err := http.ListenAndServe(host, handler)
	if err != nil {
		logger.Error("error while listening", slog.String("error", err.Error()))
	}

	return err
}
