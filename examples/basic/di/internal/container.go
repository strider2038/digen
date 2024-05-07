// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator  (built at ).
// See docs at https://github.com/strider2038/digen

package internal

import (
	"basic/app/config"
	"basic/app/domain"
	"basic/app/httphandler"
	"basic/app/usecase"
	"basic/di/internal/factories"
	"basic/di/internal/lookup"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Container struct {
	err error

	config config.Params
	logger *log.Logger
	db     *sql.DB
	server *http.Server

	params       *ParamsContainer
	api          *APIContainer
	useCases     *UseCaseContainer
	repositories *RepositoryContainer
}

func NewContainer() *Container {
	c := &Container{}
	c.params = &ParamsContainer{Container: c}
	c.api = &APIContainer{Container: c}
	c.useCases = &UseCaseContainer{Container: c}
	c.repositories = &RepositoryContainer{Container: c}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

type ParamsContainer struct {
	*Container

	serverPort     int
	serverHost     string
	requestTimeout time.Duration
}

type APIContainer struct {
	*Container

	findEntityHandler *httphandler.FindEntity
}

type UseCaseContainer struct {
	*Container

	findEntity *usecase.FindEntity
}

type RepositoryContainer struct {
	*Container

	entityRepository domain.EntityRepository
}

func (c *Container) Config(ctx context.Context) config.Params {
	return c.config
}

func (c *Container) Logger(ctx context.Context) *log.Logger {
	if c.logger == nil && c.err == nil {
		c.logger = factories.CreateLogger(ctx, c)
	}
	return c.logger
}

func (c *Container) DB(ctx context.Context) *sql.DB {
	if c.db == nil && c.err == nil {
		c.db = factories.CreateDB(ctx, c)
	}
	return c.db
}

func (c *Container) Server(ctx context.Context) *http.Server {
	if c.server == nil && c.err == nil {
		c.server = factories.CreateServer(ctx, c)
	}
	return c.server
}

func (c *Container) Params() lookup.ParamsContainer {
	return c.params
}

func (c *ParamsContainer) ServerPort(ctx context.Context) int {
	if c.serverPort == 0 && c.err == nil {
		c.serverPort = factories.CreateParamsServerPort(ctx, c)
	}
	return c.serverPort
}

func (c *ParamsContainer) ServerHost(ctx context.Context) string {
	if c.serverHost == "" && c.err == nil {
		c.serverHost = factories.CreateParamsServerHost(ctx, c)
	}
	return c.serverHost
}

func (c *ParamsContainer) RequestTimeout(ctx context.Context) time.Duration {
	if c.requestTimeout == 0 && c.err == nil {
		c.requestTimeout = factories.CreateParamsRequestTimeout(ctx, c)
	}
	return c.requestTimeout
}

func (c *Container) API() lookup.APIContainer {
	return c.api
}

func (c *APIContainer) FindEntityHandler(ctx context.Context) *httphandler.FindEntity {
	if c.findEntityHandler == nil && c.err == nil {
		c.findEntityHandler = factories.CreateAPIFindEntityHandler(ctx, c)
	}
	return c.findEntityHandler
}

func (c *Container) UseCases() lookup.UseCaseContainer {
	return c.useCases
}

func (c *UseCaseContainer) FindEntity(ctx context.Context) *usecase.FindEntity {
	if c.findEntity == nil && c.err == nil {
		c.findEntity = factories.CreateUseCasesFindEntity(ctx, c)
	}
	return c.findEntity
}

func (c *Container) Repositories() lookup.RepositoryContainer {
	return c.repositories
}

func (c *RepositoryContainer) EntityRepository(ctx context.Context) domain.EntityRepository {
	if c.entityRepository == nil && c.err == nil {
		c.entityRepository = factories.CreateRepositoriesEntityRepository(ctx, c)
	}
	return c.entityRepository
}

func (c *Container) SetConfig(s config.Params) {
	c.config = s
}

func (c *RepositoryContainer) SetEntityRepository(s domain.EntityRepository) {
	c.entityRepository = s
}

func (c *Container) Close() {
	if c.db != nil {
		c.db.Close()
	}

	if c.server != nil {
		c.server.Close()
	}
}
