//go:build !nogui
// +build !nogui

package main

import (
	"github.com/itohio/phingo/pkg/app"
	"github.com/itohio/phingo/pkg/repository"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newRootCmd(commands ...*cobra.Command) *cobra.Command {
	var repoUrl *string

	var rootCmd = &cobra.Command{
		Use:     "phingo",
		Version: version.Version,
		Short:   "Phinancial assistant app",
		Long:    ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			globalRepository, err = repository.New(*repoUrl)
			return err
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if globalRepository == nil {
				return nil
			}
			return globalRepository.Close()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New(globalRepository)
			return a.Run()
		},
	}

	repoUrl = rootCmd.PersistentFlags().StringP("repository", "r", ".phingo", "Data Repository path.")

	rootCmd.AddCommand(commands...)

	return rootCmd
}
