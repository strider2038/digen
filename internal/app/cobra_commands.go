package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMainCommand(opts *Options) *cobra.Command {
	command := &cobra.Command{
		Use:   "digen",
		Short: "Dependency Injection Container Generator",
		Long:  opts.description(),
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
			fmt.Printf(options.description())
		},
	}
}

func newInitCommand(options *Options) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		Short:         "Generates skeleton for Dependency Injection Containers",
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
		Short:         "Generates Dependency Injection Containers",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(options)
		},
	}
}
