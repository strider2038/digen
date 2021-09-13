package di

import "text/template"

type templateParameters struct {
	ContainerName string
	ServicePrefix string
	ServicePath   string
	ServiceName   string
	ServiceTitle  string
	ServiceType   string
	HasDefinition bool
	PanicOnNil    bool
}

type containerTemplateParameters struct {
	ContainerArguments       string
	ContainerArgumentSetters string
	ContainerMethods         string
}

type attachedContainerTemplateParameters struct {
	ContainerName                   string
	AttachedContainerName           string
	AttachedContainerTitle          string
	AttachedContainerDefinitionType string
}

var headingTemplate = template.Must(template.New("heading").Parse(`// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator {{.Version}} (built at {{.BuildTime}}).
// See docs at https://github.com/strider2038/

`))

var getterTemplate = template.Must(template.New("getter").Parse(`
func (c *{{.ContainerName}}) {{.ServiceTitle}}() {{.ServiceType}} {
{{ if .HasDefinition }}	if c.{{.ServiceName}} == nil {
		{{ if .PanicOnNil }}panic("missing {{.ServiceTitle}}"){{ else }}c.{{.ServiceName}} = definitions.Create{{.ServicePrefix}}{{.ServiceTitle}}(c){{ end }}
	}
{{ end }}	return c.{{.ServiceName}}
}
`))

var attachedContainerGetterTemplate = template.Must(template.New("internal container getter").Parse(`
func (c *{{.ContainerName}}) {{.AttachedContainerTitle}}() {{.AttachedContainerDefinitionType}} {
	return c.{{.AttachedContainerName}}
}
`))

var setterTemplate = template.Must(template.New("setter").Parse(`
func (c *{{.ContainerName}}) Set{{.ServiceTitle}}(s {{.ServiceType}}) {
	c.{{.ServiceName}} = s
}
`))

var closerTemplate = template.Must(template.New("closer").Parse(`
	if c.{{.ServiceName}} != nil {
		c.{{.ServiceName}}.Close()
	}
`))

var definitionTemplate = template.Must(template.New("factory").Parse(`
func Create{{.ServicePrefix}}{{.ServiceTitle}}(c Container) {{.ServiceType}} {
	panic("not implemented")
}
`))

var internalContainerTemplate = template.Must(template.New("internal container").Parse(`package internal

type Container struct {
	// err holds first initialization error
	err error

	// put the list of your services here
	// for example
	//  log *log.Logger
}

// Error returns first initialization error, which can be set via SetError in service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError set first error into container. It is used in public container to return initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err != nil {
		c.err = err
	}
}
`))

var publicContainerTemplate = template.Must(template.New("public container").Parse(`
type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer({{.ContainerArguments}}) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  &internal.Container{},
	}
{{.ContainerArgumentSetters}}
	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
{{.ContainerMethods}}
func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`))

var argumentTemplate = template.Must(template.New("argument").Parse(
	"\n\t{{.ServiceName}} {{.ServiceType}},",
))

var argumentSetterTemplate = template.Must(template.New("argument setter").Parse(
	"\n\tc.c.Set{{.ServiceTitle}}({{.ServiceName}})",
))

var publicGetterTemplate = template.Must(template.New("public getter").Parse(`
func (c *{{.ContainerName}}) {{.ServiceTitle}}() ({{.ServiceType}}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.{{.ServicePath}}{{.ServiceTitle}}()
	err := c.c.Error()
	if err != nil {
		return nil, err
	}

	return s, err
}
`))

var publicSetterTemplate = template.Must(template.New("public setter").Parse(`
func (c *{{.ContainerName}}) Set{{.ServiceTitle}}(s {{.ServiceType}}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.{{.ServicePath}}Set{{.ServiceTitle}}(s)
}
`))
