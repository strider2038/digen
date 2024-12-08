// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator rev-af34aad-dirty.
// See docs at https://github.com/strider2038/digen
package lookup

import (
	config "basic/app/config"
	domain "basic/app/domain"
	httphandler "basic/app/httphandler"
	usecase "basic/app/usecase"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Container interface {
	// SetError sets the first error into container. The error is used in the public container to return an initialization error.
	// Deprecated. Return error in factory instead.
	SetError(err error)

	Config(ctx context.Context) config.Params
	Logger(ctx context.Context) *log.Logger
	DB(ctx context.Context) *sql.DB
	Server(ctx context.Context) *http.Server

	Params() ParamsContainer
	API() APIContainer
	UseCases() UseCaseContainer
	Repositories() RepositoryContainer
}

type ParamsContainer interface {
	ServerPort(ctx context.Context) int
	ServerHost(ctx context.Context) string
	RequestTimeout(ctx context.Context) time.Duration
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
