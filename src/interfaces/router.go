package interfaces

import (
	"context"
	"net/http"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

type Router interface {
	Handle(ctx context.Context, routes []config.Route) http.Handler
}
