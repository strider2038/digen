package usecase

import (
	"context"
	"fmt"

	"basic/app/domain"
)

type FindEntity struct {
	entities domain.EntityRepository
}

func NewFindEntity(entities domain.EntityRepository) *FindEntity {
	return &FindEntity{entities: entities}
}

func (uc *FindEntity) Handle(ctx context.Context, id int) (*domain.Entity, error) {
	entity, err := uc.entities.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find entity: %w", err)
	}

	return entity, nil
}
