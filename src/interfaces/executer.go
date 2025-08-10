package interfaces

import (
	"context"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

type Executer interface {
	Execute(ctx context.Context, route config.Route) (result []byte, exitCode int, err error)
}
