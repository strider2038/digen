package console

func runInit(options *Options) error {
	config := newConfig()
	err := initConfig(config)
	if err != nil {
		return err
	}

	generator, err := newGenerator(options, config)
	if err != nil {
		return err
	}

	return generator.Initialize()
}
