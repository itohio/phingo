package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
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
		newProjectLogCmd(),
		newProjectSetCmd(),
		newProjectSetRateCmd(),
		newProjectSetDatesCmd(),
		newProjectShowCmd(),
		newProjectDeleteCmd(),
	)

	return cmd
}

func newProjectSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"set", "a", "s"},
		Version: version.Version,
		Short:   "set/add a project",
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

func newProjectSetRateCmd() *cobra.Command {
	var hourly *bool
	cmd := &cobra.Command{
		Use:     "set-rate",
		Aliases: []string{"sr"},
		Version: version.Version,
		Short:   "set a rate for the project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("Must supply project id/name, amount and denomination as arguments")
			}
			projects := globalRepository.Projects(args[0])
			if len(projects) != 1 {
				return errors.New("please provide a valid project id/name that results in a unique entry")
			}
			amount, err := strconv.ParseFloat(args[1], 32)
			if err != nil {
				return err
			}
			if *hourly {
				projects[0].Rate = &types.Project_Hourly{
					Hourly: &types.Price{
						Amount: float32(amount),
						Denom:  args[2],
					},
				}
			} else {
				projects[0].Rate = &types.Project_Total{
					Total: &types.Price{
						Amount: float32(amount),
						Denom:  args[2],
					},
				}
			}

			return globalRepository.SetProject(projects[0])
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	hourly = cmd.Flags().BoolP("hourly", "", false, "Set hourly rate (default is total per project)")

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
			if len(args) != 1 {
				return errors.New("must supply project id/name")
			}
			projects := globalRepository.Projects(args[0])
			if len(projects) != 1 {
				return errors.New("please provide a valid project id/name that results in a unique entry")
			}

			if d, err := sanitizeDateTime(*start); err == nil {
				projects[0].StartDate = d
			}
			if d, err := sanitizeDateTime(*end); err == nil {
				projects[0].EndDate = d
			}
			if *duration > time.Hour {
				projects[0].EndDate = time.Now().Add(*duration).Format(saneDateTimeLayout)
			}

			return globalRepository.SetProject(projects[0])
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
	cmd := &cobra.Command{
		Use:     "show",
		Version: version.Version,
		Short:   "show projects",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
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

const saneDateTimeLayout = "2006-01-02 15:04"

func sanitizeDateTime(val string) (string, error) {
	for _, fmt := range []string{
		saneDateTimeLayout,
		time.ANSIC,
		time.Kitchen,
		time.RFC1123,
		time.RubyDate,
	} {
		if t, err := time.Parse(fmt, val); err == nil {
			return t.Format(saneDateTimeLayout), nil
		}
	}
	return "", errors.New("invalid time format")
}

func newProjectLogCmd() *cobra.Command {
	var (
		completed *bool
		progress  *float32
		duration  *time.Duration
		started   *string
	)
	cmd := &cobra.Command{
		Use:       "log",
		ValidArgs: []string{"project Id/Name", "description"},
		Version:   version.Version,
		Short:     "log project progress",
		Long:      ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("project ID/name and a description must be provided")
			}
			description := strings.Join(args[1:], " ")
			if len(description) < 7 {
				return errors.New("the description must be longer than 7 letters")
			}
			startedSanitized, err := sanitizeDateTime(*started)
			if err != nil {
				return err
			}
			projects := globalRepository.Projects(args[0])
			if len(projects) != 1 {
				return errors.New("please provide a valid project id/name that results in a unique entry")
			}
			projects[0].Log = append(projects[0].Log, &types.Project_LogEntry{
				Start:       startedSanitized,
				Description: description,
				Duration:    int64(*duration),
				Progress:    *progress,
			})
			projects[0].Completed = *completed
			return globalRepository.SetProject(projects[0])
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
	started = cmd.Flags().StringP("started", "s", time.Now().Format(saneDateTimeLayout), "Date and time when the task started")

	cmd.MarkFlagRequired("progress")
	cmd.MarkFlagRequired("time-spent")

	return cmd
}
