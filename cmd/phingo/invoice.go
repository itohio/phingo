package main

import (
	"errors"
	"os"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newInvoiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invoice",
		Aliases: []string{"inv", "i"},
		Version: version.Version,
		Short:   "manage invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(
		newInvoiceSetCmd(),
		newInvoiceDelCmd(),
		newInvoiceShowCmd(),
	)

	return cmd
}

func newInvoiceDelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "del",
		Version: version.Version,
		Short:   "delete invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newInvoiceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"add", "new", "a"},
		Version: version.Version,
		Short:   "add/set invoice",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newInvoiceShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"list", "ls"},
		Version: version.Version,
		Short:   "show all invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := globalRepository.Read()
			if err != nil {
				return err
			}
			cfg := globalRepository.Config()
			export, err := engine.New("console", cfg)
			if err != nil {
				return err
			}
			tpl := globalRepository.Templates("invoices")
			if len(tpl) == 0 {
				return errors.New("please create a invoices.md template")
			}
			invoices := globalRepository.Invoices(args...)

			return export.ExportInvoices(os.Stdout, tpl[0], invoices, nil)
		},
	}

	return cmd
}
