# DIGEN - Dependency Injection Container Generator

Project is under development

## Installation

Download binary from releases page. Add path to binary into export path. For example, you can edit `~/.profile`.

```bash
export PATH=$PATH:$HOME/app/digen
```

## How to use

### Initialize a new container

To initialize project run command.

```bash
digen init
```

Choose work directory. After that container skeleton will be generated. Then simply write service name and its type in `Container` struct or in a separate container. After any update run `digen generate` command to generate container and definitions.

### Manage service definitions

For changing service definition behaviour you can use any of the tags.

* `set` tag - will generate setters for internal and public containers
* `close` tag - generate closer method call
* `required` tag - will generate argument for public container constructor
* `public` tag - will generate getter for public container
* `external` tag - no definition, panic if empty, force public setter

Example

```golang
type Container struct {
	err error

	configuration config.Configuration `di:"required"`
	logger        *log.Logger          `di:"required"`
	connection    *sql.Conn            `di:"external,close"`

	handler *httpadapter.GetEntityHandler `di:"public"`

	useCase    *UseCaseContainer    `di:"container"`
	repository *RepositoryContainer `di:"container"`
}

type UseCaseContainer struct {
	*Container

	findEntity *usecase.FindEntity
}

type RepositoryContainer struct {
	*Container

	entityRepository domain.EntityRepository `di:"set"`
}

func (c *Container) Error() error {
	return c.err
}

func (c *Container) SetError(err error) {
	if err != nil && c.err != nil {
		c.err = err
	}
}
```

## TODO

* [x] public container generator
    * [x] generate getters
    * [x] generate requirements
* [x] use cobra/viper
* [x] `SetError` method
* [x] generate package docs
* [x] skeleton generation command (init)
* [x] import definitions package
* [x] remove unnecessary imports
* [x] definitions generator
* [x] handle multiple containers
* [x] better console output via logger
* [x] definitions for multiple containers
* [x] unique names for separate container definitions
* [x] prompt for init (set work_dir, write first config)
* [ ] definitions updater
* [ ] ability to choose specific file for definition
* [ ] apply gofmt
* [ ] generate external services
* [ ] custom close functions
* [ ] custom container name
* [ ] write doc
* [ ] force variable name uniqueness
* [ ] don't check for nil on static structs (not so simple, may be option via tag `static`?)
