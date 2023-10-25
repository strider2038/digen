package app

func Execute(options ...OptionFunc) error {
	opts := &Options{}
	for _, setOption := range options {
		setOption(opts)
	}

	return newMainCommand(opts).Execute()
}
