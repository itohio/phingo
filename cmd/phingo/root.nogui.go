package main

import (
	"github.com/itohio/phingo/pkg/repository"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newRootCmd(commands ...*cobra.Command) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "phingo",
		Version: version.Version,
		Short:   "Phinancial assistant app",
		Long:    ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			globalRepository, err = repository.New(*flagRepositoryUrl)
			return err
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if globalRepository == nil {
				return nil
			}
			return globalRepository.Close()
		},
	}

	rootCmd.AddCommand(commands...)

	return rootCmd
}
