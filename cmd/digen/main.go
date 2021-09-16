package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/strider2038/digen/internal/console"
)

var (
	version string
	date    string
)

func main() {
	err := console.Execute(
		console.Version(version),
		console.BuildTime(date),
	)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
