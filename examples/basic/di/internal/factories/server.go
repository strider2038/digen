package factories

import (
	"context"
	"net/http"

	"basic/di/internal/lookup"
)

func CreateServer(ctx context.Context, c lookup.Container) *http.Server {
	return &http.Server{
		Handler: c.API().FindEntityHandler(ctx),
	}
}
