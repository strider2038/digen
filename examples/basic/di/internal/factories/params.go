package factories

import (
	"context"
	"time"

	"basic/di/internal/lookup"
)

func CreateParamsServerPort(ctx context.Context, c lookup.Container) int {
	return 3000
}

func CreateParamsServerHost(ctx context.Context, c lookup.Container) string {
	return "127.0.0.1"
}

func CreateParamsRequestTimeout(ctx context.Context, c lookup.Container) time.Duration {
	return time.Second
}
