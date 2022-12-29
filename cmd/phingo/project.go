package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/itohio/phingo/pkg/bi"
	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"prj", "pr", "p"},
		Version: version.Version,
		Short:   "manage projects",
		Long:    ``,
	}

	cmd.AddCommand(
		newProjectLogCmd(),
		newProjectSetCmd(),
		newProjectSetRateCmd(),
		newProjectSetClientCmd(),
		newProjectSetDatesCmd(),
		newProjectShowCmd(),
		newProjectDeleteCmd(),
	)

	return cmd
}

func newProjectSetCmd() *cobra.Command {
	var (
		hourly      *bool
		amount      *float32
		denom       *string
		account     *string
		client      *string
		name        *string
		description *string
		start       *string
		end         *string
		duration    *time.Duration
	)
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "set", "a"},
		Version: version.Version,
		Short:   "set/add a project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			accs := globalRepository.Accounts(*account)
			if len(accs) != 1 {
				return errors.New("please provide a valid account ID")
			}
			cls := globalRepository.Clients(*client)
			if len(cls) != 1 {
				return errors.New("please provide a valid client ID")
			}

			p := &types.Project{
				Name:        *name,
				Description: *description,
				Client:      cls[0],
				Account:     accs[0],
			}
			p.SetRate(*amount, *denom, *hourly)
			if d, err := bi.SanitizeDateTime(*start); err == nil {
				p.StartDate = d
			}
			if d, err := bi.SanitizeDateTime(*end); err == nil {
				p.EndDate = d
			}
			if *duration > time.Hour {
				p.EndDate = bi.Format(time.Now().Add(*duration))
			}
			err := globalRepository.SetProject(p)
			if err != nil {
				return err
			}
			log.Println("Project Id: ", p.Id)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	hourly = cmd.Flags().BoolP("hourly", "", false, "Set hourly rate (default is total per project)")
	amount = cmd.Flags().Float32P("amount", "", 0, "Provide the amount")
	account = cmd.Flags().StringP("account", "", "", "Account ID")
	client = cmd.Flags().StringP("client", "", "", "Client ID")
	denom = cmd.Flags().StringP("denom", "", "Eur", "main project denomination")
	name = cmd.Flags().StringP("name", "", "", "project short name")
	description = cmd.Flags().StringP("description", "", "Eur", "project description")
	start = cmd.Flags().StringP("start", "", "", "Set project start date")
	end = cmd.Flags().StringP("end", "", "", "Set project end date")
	duration = cmd.Flags().DurationP("duration", "d", 0, "Set project end date <duration> from now in the future")

	cmd.MarkFlagRequired("account")
	cmd.MarkFlagRequired("client")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("description")
	cmd.MarkFlagRequired("start")

	return cmd
}

func newProjectSetRateCmd() *cobra.Command {
	var (
		hourly *bool
		amount *float32
		denom  *string
	)
	cmd := &cobra.Command{
		Use:     "set-rate",
		Aliases: []string{"sr"},
		Version: version.Version,
		Short:   "set a rate for the project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects := globalRepository.Projects(args...)
			for _, p := range projects {
				p.SetRate(*amount, *denom, *hourly)
				err := globalRepository.SetProject(p)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(projects), " projects")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	hourly = cmd.Flags().BoolP("hourly", "", false, "Set hourly rate (default is total per project)")
	amount = cmd.Flags().Float32P("amount", "a", 0, "Provide the amount")
	denom = cmd.Flags().StringP("denom", "d", "Eur", "Provide denomination")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func newProjectSetClientCmd() *cobra.Command {
	var (
		account *string
		client  *string
	)
	cmd := &cobra.Command{
		Use:     "set-client",
		Aliases: []string{"set-account", "sc", "sa"},
		Version: version.Version,
		Short:   "set a client/account for the project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if *account == "" && *client == "" {
				return errors.New("either Account ID or Client ID must be provided")
			}
			var (
				acc *types.Account
				cl  *types.Client
			)
			if *account != "" {
				accs := globalRepository.Accounts(*account)
				if len(accs) != 1 {
					return errors.New("please provide a valid Account ID")
				}
				acc = accs[0]
			}
			if *client != "" {
				cls := globalRepository.Clients(*client)
				if len(cls) != 1 {
					return errors.New("please provide a valid Client ID")
				}
				cl = cls[0]
			}

			projects := globalRepository.Projects(args...)
			for _, p := range projects {
				if acc != nil {
					p.Account = acc
				}
				if cl != nil {
					p.Client = cl
				}
				err := globalRepository.SetProject(p)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(projects), " projects")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	client = cmd.Flags().StringP("client", "c", "", "Provide Client ID")
	account = cmd.Flags().StringP("account", "a", "", "Provide Account ID")

	return cmd
}

func newProjectSetDatesCmd() *cobra.Command {
	var (
		start    *string
		end      *string
		duration *time.Duration
	)
	cmd := &cobra.Command{
		Use:     "set-dates",
		Aliases: []string{"dates", "sd"},
		Version: version.Version,
		Short:   "set a start/end dates for the project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects := globalRepository.Projects(args...)

			for _, p := range projects {
				if d, err := bi.SanitizeDateTime(*start); err == nil {
					p.StartDate = d
				}
				if d, err := bi.SanitizeDateTime(*end); err == nil {
					p.EndDate = d
				}
				if *duration > time.Hour {
					p.EndDate = bi.Format(time.Now().Add(*duration))
				}
				err := globalRepository.SetProject(p)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(projects), " projects")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	start = cmd.Flags().StringP("start", "s", "", "Set project start date")
	end = cmd.Flags().StringP("end", "e", "", "Set project end date")
	duration = cmd.Flags().DurationP("duration", "d", 0, "Set project end date <duration> from now in the future")

	return cmd
}

func newProjectShowCmd() *cobra.Command {
	var (
		all   *bool
		short *bool
	)
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"list", "s", "ls"},
		Version: version.Version,
		Short:   "show projects",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects := globalRepository.Projects(args...)
			if !*all {
				for i, p := range projects {
					if p.Completed {
						projects[i] = projects[len(projects)-1]
						projects = projects[:len(projects)-1]
					}
				}
			}
			if *short {
				for _, val := range projects {
					fmt.Printf("'%s': %s", val.Name, val.Id)
					fmt.Println()
				}
				return nil
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

			return export.ExportProjects(os.Stdout, tpl[0], projects)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}
	all = cmd.Flags().BoolP("all", "a", false, "Show also completed projects")
	short = cmd.Flags().BoolP("short", "s", false, "show only names and IDs")

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

func newProjectLogCmd() *cobra.Command {
	var (
		completed   *bool
		progress    *float32
		duration    *time.Duration
		started     *string
		description *string
	)
	cmd := &cobra.Command{
		Use:       "log",
		ValidArgs: []string{"project Id/Name", "description"},
		Version:   version.Version,
		Short:     "log project progress",
		Long:      ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			startedSanitized, err := bi.SanitizeDateTime(*started)
			if err != nil {
				return err
			}
			projects := globalRepository.Projects(args[0])
			for _, p := range projects {
				p.Log = append(p.Log, &types.Project_LogEntry{
					Start:       startedSanitized,
					Description: *description,
					Duration:    int64(*duration),
					Progress:    *progress,
				})
				p.Completed = *completed
				err := globalRepository.SetProject(p)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(projects), " projects")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	completed = cmd.Flags().BoolP("completed", "c", false, "Mark the project as completed")
	progress = cmd.Flags().Float32P("progress", "p", 0, "Record the relative progress (0 = unchanged) - cumulative should add up to 100% at most")
	duration = cmd.Flags().DurationP("time-spent", "t", 0, "Time spent doing the task")
	started = cmd.Flags().StringP("started", "s", bi.Now(), "Date and time when the task started")
	description = cmd.Flags().StringP("description", "d", "", "Description")

	cmd.MarkFlagRequired("progress")
	cmd.MarkFlagRequired("description")
	cmd.MarkFlagRequired("time-spent")

	return cmd
}
