package app

import (
	"fmt"
	"strings"
)

type Options struct {
	Version   string
	BuildTime string
	DryRun    bool
}

type OptionFunc func(options *Options)

func SetVersion(version string) OptionFunc {
	return func(options *Options) {
		options.Version = strings.TrimPrefix(version, "v")
	}
}

func SetBuildTime(buildTime string) OptionFunc {
	return func(options *Options) {
		options.BuildTime = buildTime
	}
}

func (options *Options) description() string {
	buildAt := ""
	if options.BuildTime != "" {
		buildAt = " Build at " + options.BuildTime + "."
	}

	return fmt.Sprintf("DIGEN. Dependency Injection Container Generator.\nVersion %s.%s\n", options.Version, buildAt)
}
