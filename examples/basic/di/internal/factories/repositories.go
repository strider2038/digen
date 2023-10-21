package factories

import (
	"context"

	"basic/app/db"
	"basic/app/domain"
	"basic/di/internal/lookup"
)

func CreateRepositoriesEntityRepository(ctx context.Context, c lookup.Container) domain.EntityRepository {
	return db.NewEntityRepository(c.DB(ctx))
}
