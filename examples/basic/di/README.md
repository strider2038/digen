# DI container

## How to use

1. Describe your service definitions in [`definitions`](./internal/definitions) package.
2. Run `digen generate` command to regenerate container files.
3. Describe factory methods for your services in [`factories`](./internal/factories) package.
4. Build your application.

## File structure

* [`container.go`](./container.go) - generated public container
* `internal` - directory with internal packages
  * [`container.go`](./internal/container.go) - generated internal di container
  * `definitions` - package with container and service definitions (configuration file)
    * [`container.go`](./internal/definitions/container.go) - structs describing di containers (describe here your services)
  * `lookup` - directory with lookup container contracts
    * [`container.go`](./internal/lookup/container.go) - generated interfaces for internal di container (to use in factories package)
  * `factories` - package with manually written factory functions to build up services

## Service definition options

To set up service definition options use tags:

* `set` - to generate setters for internal and public containers;
* `close` - to generate closer method call;
* `required` - to generate argument for public container constructor;
* `public` - to generate getter for public container;
* `external` - no definition, panic if empty, force public setter.

## Links

See more at <https://github.com/strider2038/digen>
