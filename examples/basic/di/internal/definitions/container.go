package definitions

import (
	"database/sql"
	"log"
	"net/http"

	"basic/app/config"
	"basic/app/domain"
	"basic/app/httphandler"
	"basic/app/usecase"
)

// Container is a root dependency injection container. It is required to describe
// your services.
type Container struct {
	Config config.Params `di:"required"`
	Logger *log.Logger
	DB     *sql.DB `di:"close"`

	Server *http.Server `di:"public,close" factory-file:"server"`

	API          APIContainer
	UseCases     UseCaseContainer
	Repositories RepositoryContainer
}

type APIContainer struct {
	FindEntityHandler *httphandler.FindEntity `di:"public"`
}

type UseCaseContainer struct {
	FindEntity *usecase.FindEntity
}

type RepositoryContainer struct {
	EntityRepository domain.EntityRepository `di:"set"`
}
