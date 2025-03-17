# DIGEN - Dependency Injection Container Generator

## Installation

### Install on Linux

```shell
# binary will be $(go env GOPATH)/bin/digen
sh -c "$(curl --location https://raw.githubusercontent.com/strider2038/digen/master/install.sh)" -- -d -b $(go env GOPATH)/bin
digen version
```

### Go installer

```shell
go install github.com/strider2038/digen/cmd/digen@latest
digen version
```

## How to use

### Initialize a new container

To initialize new container skeleton run command.

```bash
digen init
```

Then describe your service definitions in the `Container` struct (`<workdir>/internal/definitions/container.go`). 
See [examples](./examples). 
After any update run `digen generate` command to generate container and factories.

### File structure

* base directory (recommended name `di`)
  * `container.go` - generated public container
  * `internal` - directory with internal packages
    * `container.go` - generated internal di container
    * `definitions` - package with container and service definitions (configuration file)
      * `container.go` - structs describing di containers (describe here your services)
    * `lookup` - directory with lookup container contracts
      * `container.go` - generated interfaces for internal di container (to use in factories package)
    * `factories` - package with manually written factory functions to build up services

### Service definition options

There are two ways to set up service definition options: by tags and by comments.
When both are present, options by tags will override options by comments (flags will be merged).

* tag `di` for flag options:
  * `set` - to generate setters for internal and public containers;
  * `close` - to generate closer method call;
  * `required` - to generate argument for public container constructor;
  * `public` - to generate getter for public container.
* tag `factory_pkg` to set up factory package;
* tag `factory_name` to set up factory filename (without extension);
* tag `public_name` to override service getter for public container.

Example of `definitions/container.go` using tags.

```golang
type Container struct {
    Configuration config.Configuration `di:"required"`
    Logger        *log.Logger          `di:"required"`
    Conn          *sql.Conn            `di:"set,close"`

    Handler *httpadapter.GetEntityHandler `di:"public"`

    UseCases     UseCaseContainer
    Repositories RepositoryContainer
}

type UseCaseContainer struct {
    FindEntity *usecase.FindEntity
}

type RepositoryContainer struct {
    EntityRepository domain.EntityRepository `di:"set"`
}
```

Example of `definitions/container.go` using comments.

```golang
type Container struct {
    // di: required
    Configuration config.Configuration
    // di: required
    Logger *log.Logger
    // di: set,close
    Conn *sql.Conn

    // di: public
    // di: factory_pkg: example.com/test/infrastructure/api/http
    // di: factory_name: handler
    Handler *httpadapter.GetEntityHandler

    UseCases     UseCaseContainer
    Repositories RepositoryContainer
}

type UseCaseContainer struct {
    FindEntity *usecase.FindEntity
}

type RepositoryContainer struct {
    // di: set
    EntityRepository domain.EntityRepository
}
```

## Configuration

DIGEN configuration can be presented in `digen.yaml`/`digen.yml`/`digen.json` file in the project root directory.

```yaml
version: v0.2
container:
  # base directory with Dependency Injection Container files
  dir: di # required
factories:
  # option can be used to disable return error by default
  returnError: true
errorHandling:
  # options for error handling
  # default values described below, can be omitted
  new:
    pkg: 'fmt'
    func: 'Errorf'
  join:
    pkg: 'errors'
    func: 'Join'
  wrap:
    pkg: 'fmt'
    func: 'Errorf'
    verb: '%w'
```

## TODO

* [x] public container generator
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
* [x] better generation with `_config.go` file
* [x] apply gofmt
* [x] move contracts into separate package
* [x] generate README.md for root package
* [x] ability to choose specific file for factory func
* [x] ability to set public definition name
* [x] check app version in config
* [x] force variable name / package name uniqueness
* [ ] custom close functions
* [ ] describe basic app example
* [ ] add complex app example with tests and fake repository
* [ ] definitions updater
* [ ] remove `SetError` method
* [ ] write doc
