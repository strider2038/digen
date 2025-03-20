# DI container

## How to use

1. Describe your service definitions in [`definitions/container.go`](./internal/definitions/container.go) file.
2. Run `digen generate` command to regenerate container files.
3. Describe factory methods for your services in [`factories`](./internal/factories) package.
4. Build your application.

## File structure

* [`container.go`](./container.go) - generated public container
* `internal` - directory with internal packages
  * [`container.go`](./internal/container.go) - generated internal di container
  * `definitions` - package with container and service definitions (configuration file)
    * [`container.go`](./internal/definitions/container.go) - structs describing di containers (describe here your services)
  * `factories` - package with manually written factory functions to build up services
* `lookup` - directory with lookup container contracts
  * [`container.go`](./lookup/container.go) - generated interfaces for internal di container (to use in factories package)

## Service definition options

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

## Links

See more at <https://github.com/strider2038/digen>
