package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client",
		Version: version.Version,
		Short:   "manage clients",
		Long:    ``,
	}

	cmd.AddCommand(
		newClientDelCmd(),
		newClientSetCmd(),
		newClientContactCmd(),
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
			for _, c := range *contact {
				kv := strings.SplitN(c, "=", 1)
				if len(kv) != 2 {
					return errors.New("contact info must be key=value")
				}
				cl.Contact[kv[0]] = kv[1]
			}
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
			if len(args) < 1 {
				return errors.New("at least one account id/name and one contact must be provided")
			}
			clients := globalRepository.Clients(args...)

			contacts := make(map[string]string, len(*contact))
			for _, c := range *contact {
				kv := strings.SplitN(c, "=", 2)
				log.Println("kv", kv)
				if len(kv) == 2 {
					contacts[kv[0]] = kv[1]
				} else {
					contacts[kv[0]] = ""
				}
			}

			for _, cl := range clients {
				if cl.Contact == nil {
					cl.Contact = make(map[string]string)
				}
				for k, v := range contacts {
					if v == "" {
						delete(cl.Contact, k)
					} else {
						cl.Contact[k] = v
					}
				}

				err := globalRepository.SetClient(cl)
				if err != nil {
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
	contact = cmd.Flags().StringArrayP("contact", "c", nil, "Key-value pair for contact information, e.g. \"Name=My name\"")

	return cmd
}

func newClientShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Version: version.Version,
		Short:   "show all clients",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := globalRepository.Config()
			export, err := engine.New("console", cfg)
			if err != nil {
				return err
			}
			tpl := globalRepository.Templates("clients")
			if len(tpl) == 0 {
				return errors.New("please create a clients.md template")
			}
			clients := globalRepository.Clients(args...)

			return export.ExportClients(os.Stdout, tpl[0], clients)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}

	return cmd
}
