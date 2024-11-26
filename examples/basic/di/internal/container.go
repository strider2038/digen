// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator rev-44eb3f4-dirty.
// See docs at https://github.com/strider2038/digen

package internal

import (
	config "basic/app/config"
	domain "basic/app/domain"
	httphandler "basic/app/httphandler"
	usecase "basic/app/usecase"
	factories "basic/di/internal/factories"
	lookup "basic/di/internal/lookup"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Container struct {
	err  error
	init bitset

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
	c.init = make(bitset, 1)
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
	if !c.init.IsSet(1) && c.err == nil {
		c.logger = factories.CreateLogger(ctx, c)
		c.init.Set(1)
	}
	return c.logger
}

func (c *Container) DB(ctx context.Context) *sql.DB {
	if !c.init.IsSet(2) && c.err == nil {
		c.db = factories.CreateDB(ctx, c)
		c.init.Set(2)
	}
	return c.db
}

func (c *Container) Server(ctx context.Context) *http.Server {
	if !c.init.IsSet(3) && c.err == nil {
		c.server = factories.CreateServer(ctx, c)
		c.init.Set(3)
	}
	return c.server
}

func (c *Container) Params() lookup.ParamsContainer {
	return c.params
}

func (c *ParamsContainer) ServerPort(ctx context.Context) int {
	if !c.init.IsSet(4) && c.err == nil {
		c.serverPort = factories.CreateParamsServerPort(ctx, c)
		c.init.Set(4)
	}
	return c.serverPort
}

func (c *ParamsContainer) ServerHost(ctx context.Context) string {
	if !c.init.IsSet(5) && c.err == nil {
		c.serverHost = factories.CreateParamsServerHost(ctx, c)
		c.init.Set(5)
	}
	return c.serverHost
}

func (c *ParamsContainer) RequestTimeout(ctx context.Context) time.Duration {
	if !c.init.IsSet(6) && c.err == nil {
		c.requestTimeout = factories.CreateParamsRequestTimeout(ctx, c)
		c.init.Set(6)
	}
	return c.requestTimeout
}

func (c *Container) API() lookup.APIContainer {
	return c.api
}

func (c *APIContainer) FindEntityHandler(ctx context.Context) *httphandler.FindEntity {
	if !c.init.IsSet(7) && c.err == nil {
		c.findEntityHandler = factories.CreateAPIFindEntityHandler(ctx, c)
		c.init.Set(7)
	}
	return c.findEntityHandler
}

func (c *Container) UseCases() lookup.UseCaseContainer {
	return c.useCases
}

func (c *UseCaseContainer) FindEntity(ctx context.Context) *usecase.FindEntity {
	if !c.init.IsSet(8) && c.err == nil {
		c.findEntity = factories.CreateUseCasesFindEntity(ctx, c)
		c.init.Set(8)
	}
	return c.findEntity
}

func (c *Container) Repositories() lookup.RepositoryContainer {
	return c.repositories
}

func (c *RepositoryContainer) EntityRepository(ctx context.Context) domain.EntityRepository {
	if !c.init.IsSet(9) && c.err == nil {
		c.entityRepository = factories.CreateRepositoriesEntityRepository(ctx, c)
		c.init.Set(9)
	}
	return c.entityRepository
}

func (c *Container) SetConfig(s config.Params) {
	c.config = s
	c.init.Set(0)
}

func (c *RepositoryContainer) SetEntityRepository(s domain.EntityRepository) {
	c.entityRepository = s
	c.init.Set(9)
}

func (c *Container) Close() {
	if c.init.IsSet(2) {
		c.db.Close()
	}
	if c.init.IsSet(3) {
		c.server.Close()
	}
}
