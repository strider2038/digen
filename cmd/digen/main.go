package main

import (
	"os"

	"github.com/carlmjohnson/versioninfo"
	"github.com/pterm/pterm"
	"github.com/strider2038/digen/internal/console"
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
	err := console.Execute(
		console.Version(version),
		console.BuildTime(date),
	)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
