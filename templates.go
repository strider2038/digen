package digen

import "text/template"

type templateParameters struct {
	ContainerName string
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

var getterTemplate = template.Must(template.New("getter").Parse(`
func (c *{{.ContainerName}}) {{.ServiceTitle}}() {{.ServiceType}} {
{{ if .HasDefinition }}	if c.{{.ServiceName}} == nil {
		{{ if .PanicOnNil }}panic("missing {{.ServiceTitle}}"){{ else }}c.{{.ServiceName}} = definitions.Create{{.ServiceTitle}}(c){{ end }}
	}
{{ end }}	return c.{{.ServiceName}}
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

	s := c.c.{{.ServiceTitle}}()
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

	c.c.Set{{.ServiceTitle}}(s)
}
`))
