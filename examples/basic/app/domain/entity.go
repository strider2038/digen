package domain

import "context"

type Entity struct {
	ID   int
	Name string
}

type EntityRepository interface {
	FindByID(ctx context.Context, id int) (*Entity, error)
}
