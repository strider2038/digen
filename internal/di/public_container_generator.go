package di

import (
	"slices"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type PublicContainerGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters

	file *FileBuilder
}

func NewPublicContainerGenerator(
	container *RootContainerDefinition,
	params GenerationParameters,
) *PublicContainerGenerator {
	file := NewFileBuilder("container.go", params.rootPackageName(), PublicPackage)

	return &PublicContainerGenerator{
		container: container,
		params:    params,
		file:      file,
	}
}

func (g *PublicContainerGenerator) Generate() (*File, error) {
	g.file.AddHeading(g.params.Version)
	g.file.AddImportAliases(g.container.Imports)

	g.file.Add(
		jen.Type().Id("Container").Struct(
			jen.Id("mu").Op("*").Qual("sync", "Mutex"),
			jen.Id("c").Op("*").Qual(g.params.packageName(InternalPackage), "Container"),
		),
		jen.Line(),
		jen.Line(),
		jen.Type().Id("Injector").
			Func().Params(jen.Id("c").Op("*").Id("Container")).Error(),
		jen.Line(),
	)

	methods := make([]jen.Code, 0, 2*len(g.container.Services))
	arguments := make([]jen.Code, 0, 1)
	argumentSetters := make([]jen.Code, 0)
	gettersCount := 0

	for _, service := range g.container.Services {
		if service.IsPublic {
			methods = append(methods, jen.Line(), jen.Line(), g.generateGetter(service, nil))
			gettersCount++
		}
		if service.HasSetter {
			methods = append(methods, jen.Line(), jen.Line(), g.generateSetter(service, nil))
		}
		if service.IsRequired {
			arguments = append(arguments, g.generateConstructorArgument(service))
			argumentSetters = append(argumentSetters, g.generateConstructorSetter(service, nil))
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			if service.IsPublic {
				methods = append(methods, jen.Line(), jen.Line(), g.generateGetter(service, attachedContainer))
				gettersCount++
			}
			if service.HasSetter {
				methods = append(methods, jen.Line(), jen.Line(), g.generateSetter(service, attachedContainer))
			}
			if service.IsRequired {
				arguments = append(arguments, g.generateConstructorArgument(service))
				argumentSetters = append(argumentSetters, g.generateConstructorSetter(service, attachedContainer))
			}
		}
	}

	arguments = append(arguments,
		jen.Id("injectors").Op("...").Id("Injector"),
	)
	if len(argumentSetters) > 0 {
		argumentSetters = append(argumentSetters, jen.Line())
	}

	g.file.Add(g.generateConstructor(arguments, argumentSetters))
	g.file.Add(methods...)
	g.file.Add(jen.Line(), g.generateCloser())
	if gettersCount > 0 {
		g.file.Add(g.generateErrorHandler()...)
	}

	return g.file.GetFile()
}

func (g *PublicContainerGenerator) generateConstructor(arguments []jen.Code, argumentSetters []jen.Code) *jen.Statement {
	return jen.Func().Id("NewContainer").
		Params(arguments...).
		Params(
			jen.Op("*").Id("Container"),
			jen.Error(),
		).
		Block(slices.Concat(
			[]jen.Code{
				jen.Id("c").Op(":=").Op("&").Id("Container").Values(jen.Dict{
					jen.Id("mu"): jen.Op("&").Qual("sync", "Mutex").Op("{}"),
					jen.Id("c"):  jen.Qual(g.params.packageName(InternalPackage), "NewContainer").Op("()"),
				}),
				jen.Line(),
			},
			argumentSetters,
			[]jen.Code{
				jen.For(
					jen.Op("_").Op(",").Id("inject").Op(":=").Range().Id("injectors").Block(
						jen.Id("err").Op(":=").Id("inject").Call(jen.Id("c")),
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Return(jen.Nil(), jen.Id("err")),
						),
					),
				),
				jen.Line(),
				jen.Return(jen.Id("c"), jen.Nil()),
			},
		)...)
}

func (g *PublicContainerGenerator) generateGetter(service *ServiceDefinition, container *ContainerDefinition) *jen.Statement {
	return jen.Func().
		Params(
			jen.Id("c").Op("*").Id("Container"),
		).
		Id(service.PublicTitle()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
		).
		Params(
			jen.Id("s").Do(g.container.Type(service.Type)),
			jen.Err().Error(),
		).
		Block(
			jen.Id("c").Dot("mu").Dot("Lock").Call(),
			jen.Defer().Id("c").Dot("mu").Dot("Unlock").Call(),
			jen.Line(),
			jen.Defer().Func().Call().Block(
				jen.If(
					jen.Id("recovered").Op(":=").Recover(),
					jen.Id("recovered").Op("!=").Nil(),
				).Block(
					jen.Err().Op("=").Id("newRecoveredError").Call(
						jen.Id("recovered"),
						jen.Id("c").Dot("c").Dot("Error").Call(),
					),
				),
			).Call(),
			jen.Line(),
			jen.Id("s").Op("=").Id("c").Dot("c").Do(g.containerPath(container)).Dot(service.Title()).Call(jen.Id("ctx")),
			jen.Id("err").Op("=").Id("c").Dot("c").Dot("Error").Call(),
			jen.Line(),
			jen.Return(jen.Id("s"), jen.Err()),
		)
}

func (g *PublicContainerGenerator) generateSetter(service *ServiceDefinition, container *ContainerDefinition) *jen.Statement {
	return jen.Func().
		Id("Set" + service.Title()).
		Params(
			jen.Id("s").Do(g.container.Type(service.Type)),
		).
		Params(
			jen.Id("Injector"),
		).
		Block(
			jen.Return(
				jen.Func().Params(jen.Id("c").Op("*").Id("Container")).Params(jen.Error()).Block(
					jen.Id("c").Dot("c").Do(g.containerPath(container)).Dot("Set"+service.Title()).Call(jen.Id("s")),
					jen.Line(),
					jen.Return(jen.Nil()),
				),
			),
		)
}

func (g *PublicContainerGenerator) generateConstructorArgument(service *ServiceDefinition) *jen.Statement {
	return jen.Id(strcase.ToLowerCamel(service.Name)).
		Do(g.container.Type(service.Type))
}

func (g *PublicContainerGenerator) generateConstructorSetter(service *ServiceDefinition, container *ContainerDefinition) *jen.Statement {
	statement := jen.Id("c").Dot("c")
	if container != nil {
		statement = statement.Do(g.containerPath(container))
	}

	return statement.Dot("Set" + service.Title()).Call(jen.Id(strcase.ToLowerCamel(service.Name)))
}

func (g *PublicContainerGenerator) generateCloser() *jen.Statement {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id("Container")).
		Id("Close").Params().
		Block(
			jen.Id("c").Dot("mu").Dot("Lock").Call(),
			jen.Defer().Id("c").Dot("mu").Dot("Unlock").Call(),
			jen.Line(),
			jen.Id("c").Dot("c").Dot("Close").Call(),
		)
}

func (g *PublicContainerGenerator) containerPath(container *ContainerDefinition) func(*jen.Statement) {
	return func(statement *jen.Statement) {
		if container == nil {
			return
		}

		statement.Dot(container.Title()).
			Op("()").
			Assert(
				jen.Op("*").Qual(g.params.packageName(InternalPackage), container.Type.Name),
			)
	}
}

func (g *PublicContainerGenerator) generateErrorHandler() []jen.Code {
	return []jen.Code{
		jen.Line(),
		jen.Func().Id("newRecoveredError").
			Params(
				jen.Id("recovered").Any(),
				jen.Id("err").Error(),
			).
			Error().
			Block(
				jen.Id("r").Op(":=").Qual(g.params.ErrorHandling.Wrap.Package, g.params.ErrorHandling.Wrap.Function).Call(
					jen.Lit("panic: %v"),
					jen.Id("recovered"),
				),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(
						g.params.joinErrors(
							jen.Id("r"),
							g.params.wrapError("previous error", jen.Err()),
						),
					),
				),
				jen.Return(jen.Id("r")),
			),
	}
}
