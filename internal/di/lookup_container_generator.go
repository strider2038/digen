package di

import (
	"github.com/dave/jennifer/jen"
)

type LookupContainerGenerator struct {
	container *RootContainerDefinition
}

func NewLookupContainerGenerator(container *RootContainerDefinition) *LookupContainerGenerator {
	return &LookupContainerGenerator{container: container}
}

func (g *LookupContainerGenerator) Generate() (*File, error) {
	file := NewFileBuilder("container.go", "lookup", LookupPackage)

	file.AddImportAliases(g.container.Imports)
	file.Add(g.generateRootContainerInterface())

	for _, attachedContainer := range g.container.Containers {
		file.Add(jen.Line())
		file.Add(g.generateContainerInterface(attachedContainer))
	}

	return file.GetFile()
}

func (g *LookupContainerGenerator) generateRootContainerInterface() *jen.Statement {
	methods := make([]jen.Code, 0, len(g.container.Services)+len(g.container.Containers)+3)
	methods = append(methods,
		jen.Commentf("SetError sets the first error into container. The error is used in the public container to return an initialization error."),
		jen.Commentf("Deprecated. Return error in factory instead."),
		jen.Id("SetError").Params(jen.Id("err").Error()),
		jen.Line(),
	)

	for _, service := range g.container.Services {
		methods = append(methods, jen.Id(service.Title()).
			Params(jen.Id("ctx").Qual("context", "Context")).
			Do(g.container.Type(service.Type)),
		)
	}

	if len(g.container.Containers) > 0 {
		methods = append(methods, jen.Line())
	}
	for _, attachedContainer := range g.container.Containers {
		methods = append(methods, jen.Id(attachedContainer.Title()).Params().Id(attachedContainer.Type.Name))
	}

	return jen.Type().Id("Container").Interface(methods...)
}

func (g *LookupContainerGenerator) generateContainerInterface(container *ContainerDefinition) *jen.Statement {
	methods := make([]jen.Code, 0, len(container.Services))

	for _, service := range container.Services {
		methods = append(methods, jen.Id(service.Title()).
			Params(jen.Id("ctx").Qual("context", "Context")).
			Do(g.container.Type(service.Type)),
		)
	}

	return jen.Type().Id(container.Type.Name).Interface(methods...)
}
