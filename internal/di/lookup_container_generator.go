package di

type LookupContainerGenerator struct {
	container *RootContainerDefinition
}

func NewLookupContainerGenerator(container *RootContainerDefinition) *LookupContainerGenerator {
	return &LookupContainerGenerator{container: container}
}

func (g *LookupContainerGenerator) Generate() (*File, error) {
	file := NewFileBuilder("container.go", "lookup", LookupPackage)

	file.AddImport(`"context"`)

	file.WriteString("\ntype Container interface {\n")
	file.WriteString("\t// SetError sets the first error into container. The error is used in the public container to return an initialization error.\n")
	file.WriteString("\tSetError(err error)\n\n")
	for _, service := range g.container.Services {
		file.AddImport(g.container.GetImport(service))
		file.WriteString("\t" + service.Title() + "(ctx context.Context) " + service.Type.String() + "\n")
	}
	if len(g.container.Containers) > 0 {
		file.WriteString("\n")
		for _, attachedContainer := range g.container.Containers {
			file.WriteString("\t" + attachedContainer.Title() + "() " + attachedContainer.Type.Name + "\n")
		}
	}
	file.WriteString("}\n")

	for _, attachedContainer := range g.container.Containers {
		file.WriteString("\ntype " + attachedContainer.Type.Name + " interface {\n")
		for _, service := range attachedContainer.Services {
			file.AddImport(g.container.GetImport(service))
			file.WriteString("\t" + service.Title() + "(ctx context.Context) " + service.Type.String() + "\n")
		}
		file.WriteString("}\n")
	}

	return file.GetFile(), nil
}
