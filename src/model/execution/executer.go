package execution

import (
	"context"

	"github.com/bdoerfchen/webcmd/src/model/process"
)

type Executer interface {
	// Run a command and return its output and status.
	Execute(ctx context.Context, config Config) (result *process.Process, exitCode int, err error)
}
