package main

import (
	"os"

	"github.com/earthboundkid/versioninfo/v2"
	"github.com/pterm/pterm"
	"github.com/strider2038/digen/internal/app"
)

var (
	version string
	date    string
)

func main() {
	if version == "" {
		version = versioninfo.Short()
		date = versioninfo.LastCommit.String()
	}
	err := app.Execute(
		app.SetVersion(version),
		app.SetBuildTime(date),
	)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
