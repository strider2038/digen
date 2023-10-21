package db

import (
	"context"
	"database/sql"
	"fmt"

	"basic/app/domain"
)

type EntityRepository struct {
	db *sql.DB
}

func NewEntityRepository(db *sql.DB) *EntityRepository {
	return &EntityRepository{db: db}
}

func (r *EntityRepository) FindByID(ctx context.Context, id int) (*domain.Entity, error) {
	return nil, fmt.Errorf("not implemented")
}
