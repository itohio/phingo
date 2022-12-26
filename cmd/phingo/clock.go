package main

import (
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newClockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clock",
		Version: version.Version,
		Short:   "manage clock in/out",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
