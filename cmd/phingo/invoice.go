package main

import (
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newInvoiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invoice",
		Version: version.Version,
		Short:   "manage invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
