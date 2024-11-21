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
	g.generateRootContainer()
	g.generateContainers()
	g.generateGetters()
	g.generateSetters()
	g.generateClosers()

	// todo: remove
	//fmt.Printf("%#v", g.file.file)

	return g.file.GetFile()
}

func (g *InternalContainerGenerator) generateRootContainer() {
	fields := make([]jen.Code, 0, len(g.container.Services)+len(g.container.Containers)+1)
	fields = append(fields, jen.Id("err").Error(), jen.Line())
	for _, service := range g.container.Services {
		fields = append(fields, jen.
			Id(strcase.ToLowerCamel(service.Name)).Do(g.container.Type(service.Type)),
		)
	}

	if len(g.container.Containers) > 0 {
		fields = append(fields, jen.Line())
	}
	constructorBlocks := make([]jen.Code, 0, 1+len(g.container.Containers))
	constructorBlocks = append(constructorBlocks,
		jen.Id("c").Op(":=").Op("&").Id("Container").Op("{}"),
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

		g.file.Add(jen.Type().Id(strings.Title(container.Type.Name)).Struct(fields...))
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
			var factoryCall jen.Code
			if service.IsExternal {
				factoryCall = jen.Panic(jen.Lit("missing " + service.Title()))
			} else {
				factoryCall = jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).Op("=").
					Qual(g.params.packageName(FactoriesPackage), "Create"+service.Title()).
					Call(jen.Id("ctx"), jen.Id("c"))
			}

			statement := jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).
				Do(service.Type.ZeroComparison()).Op("&&").
				Op("c").Dot("err").Op("==").Nil()
			block = append(block,
				jen.If(statement).Block(factoryCall),
			)
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
			Id("Set" + service.Title()).
			Params(jen.Id("s").Do(g.container.Type(service.Type))).
			Block(
				jen.Id("c").Dot(strcase.ToLowerCamel(service.Name)).Op("=").Id("s"),
			)

		g.file.Add(jen.Line(), setter)
	}
}

func (g *InternalContainerGenerator) generateClosers() {
	closers := make([]jen.Code, 0, 2)

	for _, service := range g.container.Services {
		if service.HasCloser {
			serviceName := strcase.ToLowerCamel(service.Name)
			closers = append(closers,
				jen.If(jen.Id("c").Dot(serviceName).Op("!=").Nil()).Block(
					jen.Id("c").Dot(serviceName).Dot("Close").Call(),
				),
			)
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			if service.HasCloser {
				serviceName := strcase.ToLowerCamel(service.Name)
				closers = append(closers,
					jen.If(jen.Id("c").Dot(strcase.ToLowerCamel(attachedContainer.Name)).Dot(serviceName).Op("!=").Nil()).Block(
						jen.Id("c").Dot(strcase.ToLowerCamel(attachedContainer.Name)).Dot(serviceName).Dot("Close").Call(),
					),
				)
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
