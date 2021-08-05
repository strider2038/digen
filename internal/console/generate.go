package console

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/strider2038/digen"
)

func Generate() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	container, err := digen.ParseFile(config.ContainerFilename)
	if err != nil {
		return err
	}

	fmt.Println("DI container", config.ContainerFilename, "successfully parsed")

	params := digen.DefaultGenerationParameters()

	files, err := digen.Generate(container, params)
	if err != nil {
		return err
	}

	for _, file := range files {
		dir := config.WorkDir + "/" + config.PackageDirs[file.Package]
		err = os.MkdirAll(dir, 0775)
		if err != nil {
			return errors.Wrapf(err, "failed to create dir %s", dir)
		}

		err = file.WriteTo(dir)
		if err != nil {
			return err
		}
	}

	fmt.Println("Generation completed at dir", config.WorkDir)

	return nil
}
