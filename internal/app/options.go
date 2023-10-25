package app

type Options struct {
	Version           string
	BuildTime         string
	DryRun            bool
	overrideArguments bool
}

type OptionFunc func(options *Options)

func SetVersion(version string) OptionFunc {
	return func(options *Options) {
		options.Version = version
	}
}

func SetBuildTime(buildTime string) OptionFunc {
	return func(options *Options) {
		options.BuildTime = buildTime
	}
}
