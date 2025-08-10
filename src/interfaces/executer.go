package interfaces

import (
	"context"
	"io"
)

type Executer interface {
	Execute(ctx context.Context, config ExecConfig) (result []byte, exitCode int, err error)
}

type ExecConfig struct {
	Command string            // Command file or name on PATH
	Args    []string          // Process args
	Env     map[string]string // Environment variables
	Stdin   io.Reader         // Stdin stream. Can be nil to use /dev/null
}
