// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator  (built at ).
// See docs at https://github.com/strider2038/digen

package lookup

import (
	"basic/app/config"
	"basic/app/domain"
	"basic/app/httphandler"
	"basic/app/usecase"
	"context"
	"database/sql"
	"log"
	"net/http"
)

type Container interface {
	// SetError sets the first error into container. The error is used in the public container to return an initialization error.
	SetError(err error)

	Config(ctx context.Context) config.Params
	Logger(ctx context.Context) *log.Logger
	DB(ctx context.Context) *sql.DB
	Server(ctx context.Context) *http.Server

	API() APIContainer
	UseCases() UseCaseContainer
	Repositories() RepositoryContainer
}

type APIContainer interface {
	FindEntityHandler(ctx context.Context) *httphandler.FindEntity
}

type UseCaseContainer interface {
	FindEntity(ctx context.Context) *usecase.FindEntity
}

type RepositoryContainer interface {
	EntityRepository(ctx context.Context) domain.EntityRepository
}