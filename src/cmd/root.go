package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bdoerfchen/webcmd/src/common/execution"
	"github.com/bdoerfchen/webcmd/src/common/process"
	"github.com/bdoerfchen/webcmd/src/common/router"
	"github.com/bdoerfchen/webcmd/src/logging"
	"github.com/bdoerfchen/webcmd/src/services/chirouter"
	"github.com/bdoerfchen/webcmd/src/services/configloader"
	"github.com/bdoerfchen/webcmd/src/services/executer"
	"github.com/bdoerfchen/webcmd/src/services/server"
	"github.com/bdoerfchen/webcmd/src/services/shellexecuter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webcmd",
	Short: "The simple web server that executes shell commands",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the server",
	Example: `
  webcmd --help
  webcmd run -c webcmd.conf -p 8080 
  webcmd run --method GET --route /test -- echo Hello World!
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runExec(cmd.Context())
	},
}

var flagVerbose bool
var flagDryRun bool
var flagNoColor bool

var flagServerHost string
var flagServerPort uint16
var flagRoutePattern string
var flagRouteMethod string
var flagConfigFilePath string

func Start(ctx context.Context) {
	// Persistent flags
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Set the log level from INFO to DEBUG")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable the logger's colored format")

	// Config file path flag
	runCmd.Flags().StringVarP(&flagConfigFilePath, "config-file", "c", "", "Path to the webcmd config file")
	runCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Define whether server start should be skipped")
	// Convenient flags for easy route configuration
	runCmd.Flags().StringVarP(&flagRouteMethod, "method", "m", "", "Set the HTTP method for the default route")
	runCmd.Flags().StringVarP(&flagRoutePattern, "route", "r", "", "Set the pattern for the default route")
	// Server config flags
	runCmd.Flags().Uint16VarP(&flagServerPort, "port", "p", 0, "Set the server port")
	runCmd.Flags().StringVar(&flagServerHost, "host", "", "Set the host IP address to listen to")

	rootCmd.AddCommand(runCmd)
	rootCmd.ExecuteContext(ctx)
}

func runExec(ctx context.Context) {
	// Configure logger and setup context
	logLevel := slog.LevelInfo
	if flagVerbose {
		logLevel = slog.LevelDebug
	}
	logger := logging.New(logLevel, !flagNoColor)
	setupCtx, finishSetup := context.WithTimeout(ctx, 30*time.Second)
	setupCtx = logging.AddToContext(setupCtx, logger)
	logger.Info("server start")
	fmt.Println()

	logger.Debug("beginning server setup")

	// Read config
	logger.Debug("load server configuration")
	loader := configloader.New()
	config, err := loader.Load(setupCtx, flagConfigFilePath)
	if err != nil {
		logger.Error("failed to load config file: " + err.Error())
		shutdown(logger, false)
		return
	}
	// Merge parameters from cmd flags into app config
	mergeCommandFlags(config, logger)
	logger.Debug("server configuration loaded")

	// Check configuration
	if err = checkConfig(config, logger); err != nil {
		logger.Error(err.Error())
		shutdown(logger, false)
	}

	// Setup router with executers (proc + shell)
	router := chirouter.New(setupCtx,
		config.Routes,
		&execution.ExecuterCollection{
			Proc: executer.New(), // Normal proc executer
			Shell: shellexecuter.New( // Shell executer with pool
				config.Modules.ShellPool.Size,
				process.Template{
					Command:   config.Modules.ShellPool.Path,
					Args:      config.Modules.ShellPool.Args,
					OpenStdIn: true,
				},
			),
		},
	)
	logger.Debug("using shell pool", slog.String("shell", config.Modules.ShellPool.Path), slog.Int("size", int(config.Modules.ShellPool.Size)))
	finishSetup() // Cancel setupCtx
	logger.Debug("setup finished")
	fmt.Println() // Empty log line

	// Run server
	if flagDryRun {
		logger.Info("dry-run finished")
		shutdown(logger, true)
	}

	runCtx := logging.AddToContext(ctx, logger)
	server := server.New(config.Server)
	err = server.Run(runCtx, router.Handler())
	if err != nil {
		logger.Error(err.Error())
	}

	shutdown(logger, true)
}

func shutdown(logger *slog.Logger, ok bool) {
	logger.Info("shutting down...")

	if !ok {
		os.Exit(1)
	}
	os.Exit(0)
}
