package di

import (
	"strings"
)

type FileLocator struct {
	ContainerDir string
	ModulePath   string
}

func (l *FileLocator) GetFactoryFilePath(service *ServiceDefinition, defaultFilename string) string {
	if service.FactoryFileName != "" {
		defaultFilename = service.FactoryFileName
	}
	if service.FactoryPackage == "" {
		return l.GetPackageFilePath(FactoriesPackage, defaultFilename)
	}

	return l.GetPathByPackage(service.FactoryPackage) + "/" + defaultFilename
}

func (l *FileLocator) GetContainerFilePath(filename string) string {
	return strings.Join([]string{l.ContainerDir, filename}, "/")
}

func (l *FileLocator) GetPackageFilePath(packageType PackageType, filename string) string {
	var s strings.Builder
	s.Grow(len(l.ContainerDir) + len(packageDirs[packageType]) + len(filename) + 2)

	s.WriteString(l.ContainerDir)
	s.WriteString("/")
	if packageDirs[packageType] != "" {
		s.WriteString(packageDirs[packageType])
		s.WriteString("/")
	}
	s.WriteString(filename)

	return s.String()
}

func (l *FileLocator) GetPathByPackage(factoryPackage string) string {
	return strings.TrimPrefix(factoryPackage, l.ModulePath+"/")
}
