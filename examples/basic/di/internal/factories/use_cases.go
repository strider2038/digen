package factories

import (
	"context"

	"basic/app/usecase"
	"basic/di/lookup"
)

func CreateUseCasesFindEntity(ctx context.Context, c lookup.Container) *usecase.FindEntity {
	return usecase.NewFindEntity(c.Repositories().EntityRepository(ctx))
}
