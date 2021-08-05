package console

func Execute(options ...OptionFunc) error {
	opts := &Options{}
	for _, setOption := range options {
		setOption(opts)
	}

	mainCommand := newMainCommand(opts)
	if opts.overrideArguments {
		mainCommand.SetArgs(opts.Arguments)
	}

	return mainCommand.Execute()
}
