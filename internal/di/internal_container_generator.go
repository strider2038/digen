package di

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type InternalContainerGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters

	file *FileBuilder
}

func NewInternalContainerGenerator(container *RootContainerDefinition, params GenerationParameters) *InternalContainerGenerator {
	return &InternalContainerGenerator{
		container: container,
		params:    params,
		file:      NewFileBuilder("container.go", "internal", InternalPackage),
	}
}

func (g *InternalContainerGenerator) Generate() (*File, error) {
	g.file.AddImportAliases(g.container.Imports)

	g.generateRootContainer()
	g.generateContainers()
	g.generateGetters()
	g.generateSetters()
	g.generateClosers()

	return g.file.GetFile()
}

func (g *InternalContainerGenerator) generateRootContainer() {
	fields := make([]jen.Code, 0, len(g.container.Services)+len(g.container.Containers)+3)
	fields = append(fields,
		jen.Id("err").Error(),
		jen.Id("init").Qual("", "bitset"),
		jen.Line(),
	)
	for _, service := range g.container.Services {
		fields = append(fields, jen.
			Id(strcase.ToLowerCamel(service.Name)).Do(g.container.Type(service.Type)),
		)
	}

	if len(g.container.Containers) > 0 {
		fields = append(fields, jen.Line())
	}
	constructorBlocks := make([]jen.Code, 0, 2+len(g.container.Containers))
	constructorBlocks = append(constructorBlocks,
		jen.Id("c").Op(":=").Op("&").Id("Container").Op("{}"),
		jen.Id("c").Dot("init").Op("=").Make(jen.Id("bitset"), jen.Lit(g.container.ServicesCount()/64+1)),
	)

	for _, container := range g.container.Containers {
		fields = append(fields, jen.Id(strcase.ToLowerCamel(container.Name)).Op("*").Id(container.Type.Name))
		constructorBlocks = append(constructorBlocks,
			jen.Id("c").Dot(strcase.ToLowerCamel(container.Name)).
				Op("=").Op("&").Id(container.Type.Name).Values(
				jen.Id("Container").Op(":").Id("c"),
			),
		)
	}

	g.file.Add(jen.Type().Id("Container").Struct(fields...))

	constructorBlocks = append(constructorBlocks, jen.Line(), jen.Return(jen.Id("c")))
	g.file.Add(jen.Func().
		Id("NewContainer").Params().Op("*").Id("Container").
		Block(constructorBlocks...),
	)

	g.file.Add(
		jen.Line(),
		jen.Commentf("Error returns the first initialization error, which can be set via SetError in a service definition."),
		jen.Line(),
		jen.Func().
			Params(jen.Id("c").Op("*").Id("Container")).
			Id("Error").Params().Error().
			Block(jen.Return(jen.Id("c").Dot("err"))),
		jen.Line(),
		jen.Commentf("SetError sets the first error into container. The error is used in the public container to return an initialization error."),
		jen.Line(),
		jen.Func().
			Params(jen.Id("c").Op("*").Id("Container")).
			Id("SetError").Params(jen.Err().Error()).
			Block(
				jen.If(jen.Err().Op("!=").Nil().Op("&&").Id("c").Dot("err").Op("==").Nil()).Block(
					jen.Id("c").Dot("err").Op("=").Err(),
				),
			),
	)
}

func (g *InternalContainerGenerator) generateContainers() {
	for _, container := range g.container.Containers {
		fields := make([]jen.Code, 0, len(container.Services)+2)
		fields = append(fields, jen.Op("*").Id("Container"), jen.Line())

		for _, service := range container.Services {
			fields = append(fields, jen.Id(strcase.ToLowerCamel(service.Name)).Do(g.container.Type(service.Type)))
		}

		g.file.Add(
			jen.Line(),
			jen.Type().Id(strings.Title(container.Type.Name)).Struct(fields...),
		)
	}
}

func (g *InternalContainerGenerator) generateGetters() {
	g.writeServiceGetters(g.container.Services, g.container.Name)

	for _, container := range g.container.Containers {
		g.writeContainerGetter(container)
		g.writeServiceGetters(container.Services, container.Type.Name)
	}
}

func (g *InternalContainerGenerator) writeContainerGetter(container *ContainerDefinition) {
	g.file.Add(
		jen.Line(),
		jen.Func().
			Params(jen.Id("c").Op("*").Id("Container")).
			Id(container.Title()).
			Params().
			Qual(g.params.packageName(LookupPackage), container.Type.Name).
			Block(
				jen.Return(jen.Id("c").Dot(strcase.ToLowerCamel(container.Name))),
			),
	)
}

func (g *InternalContainerGenerator) writeServiceGetters(services []*ServiceDefinition, containerName string) {
	for _, service := range services {
		block := make([]jen.Code, 0, 2)
		if !service.IsRequired {
			block = append(block, g.generateInitBlock(service))
		}

		block = append(block,
			jen.Return(jen.Id("c").Dot(strcase.ToLowerCamel(service.Name))),
		)

		getter := jen.Func().Params(jen.Id("c").Op("*").Id(strings.Title(containerName))).
			Id(service.Title()).
			Params(jen.Id("ctx").Qual("context", "Context")).
			Do(g.container.Type(service.Type)).
			Block(block...)

		g.file.Add(jen.Line(), getter)
	}
}

func (g *InternalContainerGenerator) generateInitBlock(service *ServiceDefinition) *jen.Statement {
	if service.IsExternal {
		return jen.
			If(
				jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).Op("==").Nil().
					Op("&&").Op("c").Dot("err").Op("==").Nil(),
			).
			Block(
				jen.Panic(jen.Lit("missing " + service.Title())),
			)
	}

	factoryName := strings.Title(service.Prefix) + service.Title()

	withError := true
	if factory, exists := g.container.Factories[factoryName]; exists {
		withError = factory.ReturnsError
	}

	block := make([]jen.Code, 0, 2)
	if withError {
		block = append(block,
			jen.Var().Id("err").Error(),
			jen.
				List(
					jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)),
					jen.Id("err"),
				).
				Op("=").
				Qual(g.params.packageName(FactoriesPackage), "Create"+factoryName).
				Call(jen.Id("ctx"), jen.Id("c")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Id("c").Dot("SetError").Call(
					g.params.wrapError("create "+factoryName, jen.Id("err")),
				),
			),
		)
	} else {
		block = append(block,
			jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).Op("=").
				Qual(g.params.packageName(FactoriesPackage), "Create"+factoryName).
				Call(jen.Id("ctx"), jen.Id("c")),
		)
	}

	block = append(block,
		jen.Id("c").Dot("init").Dot("Set").Call(jen.Lit(service.ID)),
	)

	return jen.If(jen.Op("!").Id("c").Dot("init").Dot("IsSet").Call(jen.Lit(service.ID)).
		Op("&&").Op("c").Dot("err").Op("==").Nil()).
		Block(block...)
}

func (g *InternalContainerGenerator) generateSetters() {
	for _, service := range g.container.Services {
		g.generateSetter(g.container.Name, service)
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			g.generateSetter(attachedContainer.Type.Name, service)
		}
	}
}

func (g *InternalContainerGenerator) generateSetter(containerName string, service *ServiceDefinition) {
	if service.HasSetter || service.IsExternal || service.IsRequired {
		setter := jen.Func().
			Params(jen.Id("c").Op("*").Id(containerName)).
			Id("Set"+service.Title()).
			Params(jen.Id("s").Do(g.container.Type(service.Type))).
			Block(
				jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).Op("=").Id("s"),
				jen.Id("c").Dot("init").Dot("Set").Call(jen.Lit(service.ID)),
			)

		g.file.Add(jen.Line(), setter)
	}
}

func (g *InternalContainerGenerator) generateClosers() {
	closers := make([]jen.Code, 0, 2)

	for _, service := range g.container.Services {
		if service.HasCloser {
			closers = append(closers, g.generateCloser(service, nil))
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			if service.HasCloser {
				closers = append(closers, g.generateCloser(service, attachedContainer))
			}
		}
	}

	g.file.Add(
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id("Container")).
			Id("Close").
			Params().
			Block(closers...),
	)
}

func (g *InternalContainerGenerator) generateCloser(service *ServiceDefinition, container *ContainerDefinition) *jen.Statement {
	block := jen.Id("c")
	if container != nil {
		block = block.Dot(strcase.ToLowerCamel(container.Name))
	}
	block = block.Dot(strcase.ToLowerCamel(service.Name)).Dot("Close").Call()

	return jen.
		If(
			jen.Id("c").Dot("init").Dot("IsSet").Call(jen.Lit(service.ID)),
		).
		Block(block)
}
