package app

import (
	"fmt"

	"github.com/muonsoft/errors"
	"github.com/spf13/cobra"
)

const descriptionTemplate = `DIGEN. Dependency Injection Container Generator.
Version %s. Build at %s.`

func newMainCommand(opts *Options) *cobra.Command {
	command := &cobra.Command{
		Use:   "digen",
		Short: "Dependency Injection Container Generator",
		Long:  fmt.Sprintf(descriptionTemplate, opts.Version, opts.BuildTime),
	}

	command.PersistentFlags().BoolVar(
		&opts.DryRun,
		"dry-run",
		false,
		`Dry run will not write any changes.`,
	)

	command.AddCommand(
		newVersionCommand(opts),
		newInitCommand(opts),
		newGenerateCommand(opts),
	)

	return command
}

func newVersionCommand(options *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("DIGEN. Dependency Injection Container Generator.\nVersion %s. Build at %s.\n", options.Version, options.BuildTime)
		},
	}
}

func newInitCommand(options *Options) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		Short:         "Generates skeleton for internal container",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(options)
		},
	}
}

func newGenerateCommand(options *Options) *cobra.Command {
	return &cobra.Command{
		Use:           "generate",
		Short:         "Generates Dependency Injection Container",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := newConfig()
			err := config.ReadInConfig()
			if err != nil {
				return errors.Errorf("read config: %w", err)
			}
			if err := checkVersion(config.GetString(configAppVersion), Version); err != nil {
				return err
			}

			return newGenerator(options, config).Generate()
		},
	}
}
