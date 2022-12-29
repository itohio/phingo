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

func newAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Version: version.Version,
		Short:   "manage accounts",
		Long:    ``,
	}

	cmd.AddCommand(
		newAccountSetCmd(),
		newAccountContactsCmd(),
		newAccountDelCmd(),
		newAccountShowCmd(),
	)

	return cmd
}

func newAccountSetCmd() *cobra.Command {
	var (
		name     *string
		denom    *string
		decimals *int32
		contact  *[]string
	)

	cmd := &cobra.Command{
		Use:     "set",
		Version: version.Version,
		Short:   "set/add accounts",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if *decimals <= 0 || *decimals > 10 {
				return errors.New("decimals must be in the range (0, 10]")
			}
			acc := &types.Account{
				Name:     *name,
				Denom:    *denom,
				Decimals: *decimals,
				Contact:  make(map[string]string, len(*contact)),
			}
			for _, c := range *contact {
				kv := strings.SplitN(c, "=", 2)
				log.Println("kv", kv)
				if len(kv) != 2 {
					return errors.New("contact info must be key=value")
				}
				acc.Contact[kv[0]] = kv[1]
			}
			err := globalRepository.SetAccount(acc)
			if err != nil {
				return err
			}
			log.Println("Account Id: ", acc.Id)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	name = cmd.Flags().StringP("name", "n", "", "Unique account name")
	denom = cmd.Flags().StringP("denom", "d", "Eur", "account primary denomination")
	decimals = cmd.Flags().Int32P("decimals", "m", 2, "Number of digits after zero")
	contact = cmd.Flags().StringArrayP("contact", "c", nil, "Key-value pair for contact information, e.g. \"Name=My name\"")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("denom")

	return cmd
}

func newAccountContactsCmd() *cobra.Command {
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
			accounts := globalRepository.Accounts(args...)

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

			for _, acc := range accounts {
				if acc.Contact == nil {
					acc.Contact = make(map[string]string)
				}
				for k, v := range contacts {
					if v == "" {
						delete(acc.Contact, k)
					} else {
						acc.Contact[k] = v
					}
				}

				err := globalRepository.SetAccount(acc)
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

func newAccountDelCmd() *cobra.Command {
	var skip *bool
	cmd := &cobra.Command{
		Use:     "del",
		Version: version.Version,
		Short:   "delete accounts",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			accounts := globalRepository.Accounts(args...)
			for _, c := range accounts {
				log.Println("Deleting account id=", c.Id, "name=", c.Name)
				if err := globalRepository.DelAccount(c); err != nil {
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

func newAccountShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Version: version.Version,
		Short:   "show all accounts",
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
			tpl := globalRepository.Templates("accounts")
			if len(tpl) == 0 {
				return errors.New("please create a accounts.md template")
			}
			accounts := globalRepository.Accounts(args...)

			return export.ExportAccounts(os.Stdout, tpl[0], accounts)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}

	return cmd
}
