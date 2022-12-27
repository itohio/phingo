package main

import (
	"errors"
	"log"
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
		newProjectSetCmd(),
		newProjectShowCmd(),
		newProjectDeleteCmd(),
	)

	return cmd
}

func newProjectSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Version: version.Version,
		Short:   "add a project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
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
			projects := globalRepository.Projects(args...)

			return export.ExportProjects(os.Stdout, tpl[0], projects)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}

	return cmd
}

func newProjectDeleteCmd() *cobra.Command {
	var skip *bool
	cmd := &cobra.Command{
		Use:     "del",
		Version: version.Version,
		Short:   "delete project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects := globalRepository.Projects(args...)
			for _, c := range projects {
				log.Println("Deleting project id=", c.Id, "name=", c.Name)
				if err := globalRepository.DelProject(c); err != nil {
					if *skip {
						log.Println("Failed deleting id=", c.Id, "name=", c.Name, "err=", err.Error())
						continue
					}
					return err
				}
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	skip = cmd.Flags().BoolP("ignore-errors", "i", false, "Skip any errors")

	return cmd
}
