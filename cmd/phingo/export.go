package main

import (
	"os"

	"github.com/itohio/phingo/pkg/engine"
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

	return cmd
}
