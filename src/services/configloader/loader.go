package configloader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/bdoerfchen/webcmd/src/logging"
	"github.com/bdoerfchen/webcmd/src/model/config"
	"github.com/jinzhu/copier"
)

const DefaultConfigFile = "server.conf"

type configLoader struct{}

func New() *configLoader {
	return &configLoader{}
}

func (l *configLoader) Load(ctx context.Context, path string) (*config.AppConfig, error) {
	var useDefaultPath bool
	logger := logging.FromContext(ctx)
	result := config.DefaultAppConfig()

	// Handle no path defined
	if path == "" {
		path = DefaultConfigFile
		logger.Debug("using default config path", slog.String("path", path))
		useDefaultPath = true
	}

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") || !useDefaultPath {
			return nil, err
		}

		// Return default if no specific config file should be loaded
		return &result, nil
	}

	// Try all unmarshallers
	for _, u := range unmarshallers {
		logger.Debug("trying format '" + u.Name() + "'")

		output, err := u.Unmarshal(content)
		if err != nil {
			if errors.Is(err, errTypeNotCorrect) {
				continue
			}

			return nil, fmt.Errorf("could not read config file: %w", err)
		}

		// Go over read config and validate
		var result config.AppConfig
		result.Server = output.Server
		for _, configRoute := range output.Routes {
			copiedRoute := config.DefaultRoute()
			copier.CopyWithOption(&copiedRoute, &configRoute, copier.Option{IgnoreEmpty: true})
			copiedRoute.StatusCodes = append(config.DefaultRoute().StatusCodes, configRoute.StatusCodes...)
			result.Routes = append(result.Routes, copiedRoute)
		}

		return &result, nil
	}

	return nil, errFormatUnknown
}

func AppendRoute(appConfig *config.AppConfig, route *config.Route, logger *slog.Logger) {

}
