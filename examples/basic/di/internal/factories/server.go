package factories

import (
	"context"
	"net/http"

	"basic/di/lookup"
)

func CreateServer(ctx context.Context, c lookup.Container) *http.Server {
	return &http.Server{
		Handler:     c.API().FindEntityHandler(ctx),
		IdleTimeout: c.Params().RequestTimeout(ctx),
	}
}
