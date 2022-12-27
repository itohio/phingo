package main

import (
	"errors"
	"os"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		Version: version.Version,
		Short:   "Export phinancial records",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open("templates/default.md")
			if err != nil {
				return err
			}
			defer f.Close()
			w, err := os.Create("output.pdf")
			if err != nil {
				return err
			}
			defer w.Close()

			cfg := globalRepository.Config()

			md, err := engine.New("pdf", cfg)
			if err != nil {
				return err
			}

			_ = md

			return nil
		},
	}

	cmd.AddCommand(
		newExportProjectCmd(),
		newExportInvoiceCmd(),
	)

	return cmd
}

func newExportProjectCmd() *cobra.Command {
	var (
		how      *string
		template *string
		output   *string
		cmd      = &cobra.Command{
			Use:     "export",
			Version: version.Version,
			Short:   "Export phinancial records",
			Long:    ``,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return globalRepository.Read()
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				eng, err := engine.New(*how, globalRepository.Config())
				if err != nil {
					return err
				}

				tpl := globalRepository.Templates(*template)
				if len(tpl) == 0 {
					return errors.New("could not find a template")
				}
				if len(tpl) > 1 {
					return errors.New("there are multiple templates with such name")
				}

				projects := globalRepository.Projects(args...)

				w, err := os.Create(*output)
				if err != nil {
					return err
				}
				defer w.Close()

				return eng.ExportProjects(w, tpl[0], projects)
			},
		}
	)
	how = cmd.Flags().StringP("how", "w", "pdf", "Defines how to export the entries (possible values are html and pdf)")
	template = cmd.Flags().StringP("template", "t", "projects", "what template to use")
	output = cmd.Flags().StringP("output", "o", "projects.pdf", "path to file to write output to")
	return cmd
}

func newExportInvoiceCmd() *cobra.Command {
	var (
		how      *string
		template *string
		output   *string
		account  *string
		cmd      = &cobra.Command{
			Use:     "export",
			Version: version.Version,
			Short:   "Export phinancial records",
			Long:    ``,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return globalRepository.Read()
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				eng, err := engine.New(*how, globalRepository.Config())
				if err != nil {
					return err
				}

				tpl := globalRepository.Templates(*template)
				if len(tpl) == 0 {
					return errors.New("could not find a template")
				}
				if len(tpl) > 1 {
					return errors.New("there are multiple templates with such name")
				}

				var acc *types.Account
				accs := globalRepository.Accounts(*account)
				if len(accs) == 1 {
					acc = accs[0]
				}
				if len(accs) > 1 {
					return errors.New("there are multiple accounts with such name/id")
				}
				invoices := globalRepository.Invoices(args...)

				w, err := os.Create(*output)
				if err != nil {
					return err
				}
				defer w.Close()

				return eng.ExportInvoices(w, tpl[0], invoices, acc)
			},
		}
	)
	how = cmd.Flags().StringP("how", "w", "pdf", "Defines how to export the entries (possible values are html and pdf)")
	template = cmd.Flags().StringP("template", "t", "projects", "what template to use")
	output = cmd.Flags().StringP("output", "o", "projects.pdf", "path to file to write output to")
	account = cmd.Flags().StringP("account", "a", "", "account id/name to use for the invoice instead of the one specified in the project")
	return cmd
}
