package console

func Execute(options ...OptionFunc) error {
	opts := &Options{}
	for _, setOption := range options {
		setOption(opts)
	}

	command := newMainCommand(opts)
	if opts.overrideArguments {
		command.SetArgs(opts.Arguments)
	}

	return command.Execute()
}
