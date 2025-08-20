package execution

import (
	"fmt"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

type ExecuterCollection struct {
	Proc  Executer
	Shell Executer
}

func (c *ExecuterCollection) For(route *config.Route) (Executer, error) {
	switch {
	case route.Exec.Proc != nil:
		return c.Proc, nil
	case route.Exec.Shell != nil:
		return c.Shell, nil
	}

	return nil, fmt.Errorf("route without execution config")
}
