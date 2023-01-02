package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		Aliases: []string{"exp", "ex", "e"},
		Version: version.Version,
		Short:   "Export phinancial records",
		Long:    ``,
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
			Use:     "project",
			Aliases: []string{"prj", "pr", "p"},
			Version: version.Version,
			Short:   "Export projects records",
			Long:    ``,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return globalRepository.Read()
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				eng, err := engine.New(*how, globalRepository.Config(), globalRepository.FS())
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

				fname := *output
				if !strings.Contains(*output, ".") {
					fname = fmt.Sprintf("%s.%s", *output, *how)
				}

				log.Println("Exporting to ", *output)
				w, err := os.Create(fname)
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
	output = cmd.Flags().StringP("output", "o", "projects", "path to file to write output to")
	return cmd
}

func newExportInvoiceCmd() *cobra.Command {
	var (
		how      *string
		template *string
		output   *string
		account  *string
		cmd      = &cobra.Command{
			Use:     "invoice",
			Aliases: []string{"inv", "i"},
			Version: version.Version,
			Short:   "Export invoice records",
			Long:    ``,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return globalRepository.Read()
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				config := globalRepository.Config()
				eng, err := engine.New(*how, config, globalRepository.FS())
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
				if *account != "" {
					accs := globalRepository.Accounts(*account)
					if len(accs) == 1 {
						acc = accs[0]
					}
					if len(accs) > 1 {
						return errors.New("there are multiple accounts with such name/id")
					}
				}

				invoices := globalRepository.Invoices(args...)
				if len(invoices) == 0 {
					log.Println("Nothing to do")
					return nil
				}
				if *output == "" {
					*output = fmt.Sprintf("%s.%s", invoices[0].MakeFileName(), *how)
				}

				log.Println("Exporting to ", *output)
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
	output = cmd.Flags().StringP("output", "o", "", "path to file to write output to defaults to a file name according to invoice naming format)")
	account = cmd.Flags().StringP("account", "a", "", "account id/name to use for the invoice instead of the one specified in the project")
	return cmd
}
