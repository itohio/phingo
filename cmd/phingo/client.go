package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cl", "c"},
		Version: version.Version,
		Short:   "manage clients",
		Long:    ``,
	}

	cmd.AddCommand(
		newClientSetCmd(),
		newClientContactCmd(),
		newClientNoteCmd(),
		newClientDelCmd(),
		newClientShowCmd(),
	)

	return cmd
}

func newClientDelCmd() *cobra.Command {
	var skip *bool
	cmd := &cobra.Command{
		Use:     "del",
		Version: version.Version,
		Short:   "delete clients",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			clients := globalRepository.Clients(args...)
			for _, c := range clients {
				log.Println("Deleting client id=", c.Id, "name=", c.Name)
				if err := globalRepository.DelClient(c); err != nil {
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

func newClientSetCmd() *cobra.Command {
	var (
		name        *string
		description *string
		contact     *[]string
		notes       *[]string
	)

	cmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"new", "set", "a"},
		Version: version.Version,
		Short:   "set/add clients",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := &types.Client{
				Name:        *name,
				Description: *description,
				Notes:       *notes,
				Contact:     make(map[string]string, len(*contact)),
			}
			parseKeyValue(cl.Contact, *contact)
			err := globalRepository.SetClient(cl)
			if err != nil {
				return err
			}
			log.Println("Client Id: ", cl.Id)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	name = cmd.Flags().StringP("name", "n", "", "Unique client name")
	description = cmd.Flags().StringP("denom", "d", "Eur", "client description")
	contact = cmd.Flags().StringArrayP("contact", "c", nil, "Key-value pair for contact information, e.g. \"Name=My name\"")
	notes = cmd.Flags().StringArrayP("note", "t", nil, "A list of notes")
	cmd.MarkFlagRequired("name")

	return cmd

}

func newClientContactCmd() *cobra.Command {
	var (
		contact *[]string
	)

	cmd := &cobra.Command{
		Use:     "contact",
		Version: version.Version,
		Short:   "set/add/delete contacts",
		Long:    `will delete contact key if value is empty`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clients := globalRepository.Clients(args...)

			for _, cl := range clients {
				if cl.Contact == nil {
					cl.Contact = make(map[string]string)
				}
				parseKeyValue(cl.Contact, *contact)
				err := globalRepository.SetClient(cl)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(clients), " clients")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	contact = cmd.Flags().StringArrayP("contact", "c", nil, "Key-value pair for contact information, e.g. \"Name=My name\"")

	return cmd
}

func newClientNoteCmd() *cobra.Command {
	var (
		note *[]string
	)

	cmd := &cobra.Command{
		Use:     "note",
		Version: version.Version,
		Short:   "add notes",
		Long:    `will add notes to the client`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clients := globalRepository.Clients(args...)

			for _, cl := range clients {
				cl.Notes = append(cl.Notes, *note...)
				err := globalRepository.SetClient(cl)
				if err != nil {
					return err
				}
			}
			log.Println("Modified ", len(clients), " clients")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	note = cmd.Flags().StringArrayP("note", "t", nil, "A note to add to the client notes")

	return cmd
}

func newClientShowCmd() *cobra.Command {
	var short *bool
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"list", "s", "ls"},
		Version: version.Version,
		Short:   "show all clients",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			clients := globalRepository.Clients(args...)
			if *short {
				for _, val := range clients {
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
			tpl := globalRepository.Templates("clients")
			if len(tpl) == 0 {
				return errors.New("please create a clients.md template")
			}

			return export.ExportClients(os.Stdout, tpl[0], clients)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}
	short = cmd.Flags().BoolP("short", "s", false, "show only names and IDs")

	return cmd
}
