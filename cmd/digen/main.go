package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/strider2038/digen"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("filename is required")
		os.Exit(1)
	}

	filename := os.Args[1]

	container, err := digen.ParseFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Container", filename, "successfully parsed")

	path := filepath.Dir(filename)
	_, err = digen.Generate(container, digen.GenerationParameters{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Generation completed at path", path)
}
