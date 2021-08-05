package main

import (
	"fmt"
	"os"

	"github.com/strider2038/digen/internal/console"
)

var (
	version   string
	buildTime string
)

func main() {
	err := console.Execute(
		console.Version(version),
		console.BuildTime(buildTime),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
