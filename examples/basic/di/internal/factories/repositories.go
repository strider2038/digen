package factories

import (
	"context"

	"basic/app/db"
	"basic/app/domain"
	"basic/di/lookup"
)

func CreateRepositoriesEntityRepository(ctx context.Context, c lookup.Container) domain.EntityRepository {
	return db.NewEntityRepository(c.DB(ctx))
}
