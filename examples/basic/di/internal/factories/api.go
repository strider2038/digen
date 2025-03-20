package factories

import (
	"context"

	"basic/app/httphandler"
	"basic/di/lookup"
)

func CreateAPIFindEntityHandler(ctx context.Context, c lookup.Container) *httphandler.FindEntity {
	return httphandler.NewFindEntity(c.UseCases().FindEntity(ctx))
}
