package config

import "github.com/bdoerfchen/webcmd/src/services/server"

type AppConfig struct {
	Server server.Config
	Routes []Route
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		Server: server.Config{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Routes: make([]Route, 0),
	}
}
