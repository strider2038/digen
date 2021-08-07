package console

import (
	"fmt"

	"github.com/strider2038/digen"
)

func Init() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	params := digen.GenerationParameters{}
	file, err := digen.GenerateContainer(params)
	if err != nil {
		return err
	}

	writer := digen.NewWriter(config.WorkDir)

	err = writer.WriteFile(file)
	if err != nil {
		return err
	}

	fmt.Println("Init completed: file", file.Name, "generated")

	return nil
}
