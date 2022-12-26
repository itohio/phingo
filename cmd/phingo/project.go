package main

import (
	"errors"
	"os"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Version: version.Version,
		Short:   "manage projects",
		Long:    ``,
	}

	cmd.AddCommand(
		newProjectAddCmd(),
		newProjectShowCmd(),
		newProjectUpdateCmd(),
		newProjectDeleteCmd(),
		newProjectExportCmd(),
	)

	return cmd
}

func newProjectAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Version: version.Version,
		Short:   "add a project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newProjectShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Version: version.Version,
		Short:   "show projects",
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
			tpl := globalRepository.Templates("projects")
			if len(tpl) == 0 {
				return errors.New("please create a projects.md template")
			}
			acc := globalRepository.Accounts()
			if len(acc) == 0 {
				return errors.New("please add an account")
			}
			projects := globalRepository.Projects(args...)

			return export.ExportProjects(os.Stdout, tpl[0], projects, acc[0])
		},
	}

	return cmd
}

func newProjectDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Version: version.Version,
		Short:   "delete projects",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newProjectUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Version: version.Version,
		Short:   "update a project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func newProjectExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		Version: version.Version,
		Short:   "export projects",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
