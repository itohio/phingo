package main

import (
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Version: version.Version,
		Short:   "Start Phingo backend",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
