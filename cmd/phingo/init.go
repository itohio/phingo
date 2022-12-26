package main

import (
	"errors"

	"github.com/itohio/phingo/pkg/repository"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Version: version.Version,
		Short:   "Initialize phinancial records datastore",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("path must be provided")
			}

			return repository.Init(args[0])
		},
	}

	return cmd
}
