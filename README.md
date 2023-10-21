# DIGEN - Dependency Injection Container Generator

## Installation

### Go installer

```shell
go install github.com/strider2038/digen@latest
```

### Install on linux

```shell
sh -c "$(curl --location https://raw.githubusercontent.com/strider2038/digen/master/install.sh)" -- -d -b /usr/local/bin
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

To set up service definition options use tags:

* `set` - to generate setters for internal and public containers;
* `close` - to generate closer method call;
* `required` - to generate argument for public container constructor;
* `public` - to generate getter for public container;
* `external` - no definition, panic if empty, force public setter.

Example of `definitions/container.go`

```golang
type Container struct {
    Configuration config.Configuration `di:"required"`
    Logger        *log.Logger          `di:"required"`
    Conn          *sql.Conn            `di:"external,close"`

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

## Known issues

### Struct field is compared with nil for internal container

Workaround: don't use struct by values as services or set `required` option to generate constructor.

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
* [ ] custom close functions
* [ ] add complex app example with tests and fake repository
* [ ] ability to set definition name
* [ ] parse from multiple definition files (may encounter potential conflicts for imports)
* [ ] definitions updater
* [ ] write doc
* [ ] force variable name / package name uniqueness
