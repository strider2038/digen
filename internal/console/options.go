package console

type Options struct {
	Version           string
	BuildTime         string
	Arguments         []string
	DryRun            bool
	overrideArguments bool
}

type OptionFunc func(options *Options)

func Version(version string) OptionFunc {
	return func(options *Options) {
		options.Version = version
	}
}

func BuildTime(buildTime string) OptionFunc {
	return func(options *Options) {
		options.BuildTime = buildTime
	}
}

func Arguments(args []string) OptionFunc {
	return func(options *Options) {
		options.Arguments = args
		options.overrideArguments = true
	}
}
