package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/itohio/phingo/pkg/engine"
	"github.com/itohio/phingo/pkg/types"
	"github.com/itohio/phingo/pkg/version"
	"github.com/spf13/cobra"
)

func newInvoiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invoice",
		Aliases: []string{"inv", "i"},
		Version: version.Version,
		Short:   "manage invoices",
		Long:    ``,
	}

	cmd.AddCommand(
		newInvoiceSetCmd(),
		newInvoiceItemCmd(),
		newInvoicePaymentCmd(),
		newInvoiceDelCmd(),
		newInvoiceShowCmd(),
	)

	return cmd
}

func newInvoiceDelCmd() *cobra.Command {
	var skip *bool
	cmd := &cobra.Command{
		Use:     "del",
		Version: version.Version,
		Short:   "delete invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			invoices := globalRepository.Invoices(args...)
			for _, inv := range invoices {
				log.Println("Deleting invoice id=", inv.Id, "project=", inv.Project.Name)
				if err := globalRepository.DelInvoice(inv); err != nil {
					if *skip {
						log.Println("Failed deleting id=", inv.Id, "project=", inv.Project.Name, "err=", err.Error())
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

func newInvoiceSetCmd() *cobra.Command {
	var (
		project     *string
		client      *string
		account     *string
		code        *string
		issueDate   *string
		dueDate     *string
		marketValue *float32
		marketDenom *string
	)
	cmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"add", "new", "a"},
		Version: version.Version,
		Short:   "add/set invoice",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				prj *types.Project
				acc *types.Account
				cl  *types.Client
				mv  *types.Price
			)
			prjs := globalRepository.Projects(*project)
			if len(prjs) == 1 {
				prj = prjs[0]
				acc = prj.Account
				cl = prj.Client
			}
			accs := globalRepository.Accounts(*account)
			if len(accs) == 1 {
				acc = accs[0]
			}
			cls := globalRepository.Clients(*client)
			if len(cls) == 1 {
				cl = cls[0]
			}
			if acc == nil || cl == nil {
				return errors.New("account and client must be either provided or derived from a project")
			}

			now := time.Now()
			duePeriod := time.Duration(acc.InvoiceDuePeriod)
			if duePeriod == 0 {
				duePeriod = 3 * 24
			}
			if *issueDate != "" {
				t, err := types.ParseTime(*issueDate)
				if err != nil {
					return err
				}
				now = t
			}
			due := now.Add(time.Hour * duePeriod)
			if *dueDate != "" {
				t, err := types.ParseTime(*dueDate)
				if err != nil {
					return err
				}
				due = t
			}
			*issueDate = types.FormatTime(now)
			*dueDate = types.FormatTime(due)
			if due.Sub(now) < time.Hour {
				return errors.New("due date must be at least 1 hour in the future")
			}

			if *code == "" {
				totalInvoicesCount := globalRepository.InvoicesCount()
				invoices := globalRepository.Invoices(fmt.Sprintf("year:%d", now.Year()))
				*code = acc.MakeInvoiceCode(invoices, totalInvoicesCount, now)
			}
			if *marketValue > 0 && *marketDenom != "" {
				mv = &types.Price{
					Denom:  *marketDenom,
					Amount: *marketValue,
				}
			}

			inv := &types.Invoice{
				Project:     prj,
				Account:     acc,
				Client:      cl,
				Code:        *code,
				IssueDate:   *issueDate,
				DueDate:     *dueDate,
				MarketValue: mv,
			}
			inv.Id = inv.MakeId()
			err := globalRepository.SetInvoice(inv)
			if err != nil {
				return err
			}

			log.Println("Invoice Id: ", inv.Id)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	project = cmd.Flags().StringP("project", "p", "", "Associated Project")
	account = cmd.Flags().StringP("account", "a", "", "Associated Account (override project's account)")
	client = cmd.Flags().StringP("client", "c", "", "Associated Client (override project's client)")
	code = cmd.Flags().StringP("code", "o", "", "Invoice code (default: auto generated)")
	issueDate = cmd.Flags().StringP("issue-date", "i", "", "Invoice issue date (default: auto generated)")
	dueDate = cmd.Flags().StringP("due-date", "u", "", "Invoice due date (default: auto generated - +30 days)")
	marketValue = cmd.Flags().Float32P("market-value", "v", 0, "new currency market value with respect account/client currency (e.g. BTC value in Eur)")
	marketDenom = cmd.Flags().StringP("denom", "m", "", "Override total calculations using a different currency (w.r.t. account currency)")

	cmd.MarkFlagsRequiredTogether("market-value", "denom")

	return cmd
}

func newInvoiceItemCmd() *cobra.Command {
	var (
		name        *string
		unit        *string
		amount      *float32
		rate        *float32
		extra       *bool
		del         *bool
		index       *int
		projectItem *int
	)
	cmd := &cobra.Command{
		Use:     "item",
		Aliases: []string{"it", "i"},
		Version: version.Version,
		Short:   "add/set/remove invoice items",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			invoices := globalRepository.Invoices(args...)
			if len(invoices) != 1 {
				return errors.New("a single invoice must be specified")
			}
			inv := invoices[0]
			if *index < 0 {
				*index = len(inv.Items) + *index
			}
			if (cmd.Flags().Lookup("index").Changed || *del) && (len(inv.Items) <= *index || *index < 0) {
				return errors.New("item index out of range")
			}
			if *del {
				items, err := types.Remove(inv.Items, *index)
				if err != nil {
					return err
				}
				inv.Items = items
				return globalRepository.SetInvoice(inv)
			}

			if !cmd.Flags().Lookup("project-item").Changed {
				if inv.Project == nil {
					return errors.New("project is not set")
				}
				if *projectItem < 0 {
					*projectItem = len(inv.Project.Log) + *projectItem
				}
				if len(inv.Project.Log) <= *projectItem {
					item := inv.Project.Log[*projectItem]
					*name = item.Description
					switch r := inv.Project.Rate.(type) {
					case *types.Project_Hourly:
						*amount = float32(time.Duration(item.Duration).Hours())
						*rate = r.Hourly
						*unit = "Hours"
					case *types.Project_Total:
						*amount = item.Progress
						*rate = r.Total
						*unit = "Pcs"
					}
				}
			}

			if !cmd.Flags().Lookup("index").Changed {
				if *name == "" {
					return errors.New("name must be specified")
				}
				inv.Items = append(
					inv.Items,
					&types.Invoice_Item{
						Name:   *name,
						Amount: *amount,
						Extra:  *extra,
						Unit:   *unit,
						Rate:   *rate,
					},
				)
			} else {
				item := inv.Items[*index]
				if cmd.Flags().Lookup("name").Changed {
					item.Name = *name
				}
				if cmd.Flags().Lookup("amount").Changed {
					item.Amount = *amount
				}
				if cmd.Flags().Lookup("unit").Changed {
					item.Unit = *unit
				}
				if cmd.Flags().Lookup("rate").Changed {
					item.Rate = *rate
				}
			}
			return globalRepository.SetInvoice(inv)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}

	name = cmd.Flags().StringP("name", "n", "", "Name of the billable item")
	unit = cmd.Flags().StringP("unit", "u", "Hours", "Unit of the billable item")
	amount = cmd.Flags().Float32P("amount", "a", 0, "Amount")
	rate = cmd.Flags().Float32P("rate", "t", 0, "Rate per item")
	extra = cmd.Flags().BoolP("extra", "e", false, "If true, then this item will be applied after subtotal calculations")
	del = cmd.Flags().BoolP("delete", "D", false, "If true, deletes the item from invoice")
	index = cmd.Flags().IntP("index", "i", 0, "Index of the invoice item to modify (negative numbers start from the end)")
	projectItem = cmd.Flags().IntP("project-item", "p", 0, "Import item from project items list using the index (negative numbers start from the end)")

	cmd.MarkFlagsMutuallyExclusive("project-item", "name")
	cmd.MarkFlagsMutuallyExclusive("project-item", "amount")
	cmd.MarkFlagsMutuallyExclusive("project-item", "delete")
	cmd.MarkFlagsMutuallyExclusive("project-item", "rate")
	return cmd
}

func newInvoicePaymentCmd() *cobra.Command {
	var (
		name    *string
		date    *string
		comment *string
		amount  *float32
		denom   *string
		del     *bool
		index   *int
	)
	cmd := &cobra.Command{
		Use:     "payment",
		Aliases: []string{"payed", "pay"},
		Version: version.Version,
		Short:   "add/set/remove invoice payments",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			invoices := globalRepository.Invoices(args...)
			if len(invoices) != 1 {
				return errors.New("a single invoice must be specified")
			}
			inv := invoices[0]
			if *index < 0 {
				*index = len(inv.Payments) + *index
			}
			if (cmd.Flags().Lookup("index").Changed || *del) && (len(inv.Payments) <= *index || *index < 0) {
				return errors.New("payment index out of range")
			}
			if *del {
				payments, err := types.Remove(inv.Payments, *index)
				if err != nil {
					return err
				}
				inv.Payments = payments
				return globalRepository.SetInvoice(inv)
			}
			if cmd.Flags().Lookup("date").Changed {
				var err error
				*date, err = types.SanitizeDateTime(*date)
				if err != nil {
					return err
				}
			} else {
				*date = types.Now()
			}
			if !cmd.Flags().Lookup("index").Changed {
				inv.Payments = append(
					inv.Payments,
					&types.Invoice_Payment{
						Name: *name,
						Amount: &types.Price{
							Amount: *amount,
							Denom:  *denom,
						},
						Comment: *comment,
						Date:    *date,
					},
				)
			} else {
				payment := inv.Payments[*index]
				if cmd.Flags().Lookup("name").Changed {
					payment.Name = *name
				}
				if cmd.Flags().Lookup("date").Changed {
					payment.Date = *date
				}
				if cmd.Flags().Lookup("comment").Changed {
					payment.Comment = *comment
				}
				if cmd.Flags().Lookup("amount").Changed {
					payment.Amount = &types.Price{
						Amount: *amount,
						Denom:  *denom,
					}
				}
			}
			return globalRepository.SetInvoice(inv)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Write()
		},
	}
	name = cmd.Flags().StringP("name", "n", "", "Name of the billable item")
	comment = cmd.Flags().StringP("comment", "c", "", "Name of the billable item")
	date = cmd.Flags().StringP("date", "d", "", "Date when the payment was received")
	del = cmd.Flags().BoolP("delete", "D", false, "If true, deletes the item from invoice")
	index = cmd.Flags().IntP("index", "i", 0, "Index of the invoice item to modify (negative numbers start from the end)")
	amount = cmd.Flags().Float32P("amount", "a", 0, "Amount")
	cmd.MarkFlagsRequiredTogether("delete", "index")

	return cmd
}

func newInvoiceShowCmd() *cobra.Command {
	var short *bool
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"list", "ls"},
		Version: version.Version,
		Short:   "show all invoices",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			invoices := globalRepository.Invoices(args...)
			if *short {
				for _, val := range invoices {
					fmt.Printf("%d: '%s': %s", val.Year(), val.Id, val.Code)
					fmt.Println()
				}
				return nil
			}
			cfg := globalRepository.Config()
			export, err := engine.New("console", cfg, globalRepository.FS())
			if err != nil {
				return err
			}
			tpl := globalRepository.Templates("invoices")
			if len(tpl) == 0 {
				return errors.New("please create a invoices.md template")
			}

			return export.ExportInvoices(os.Stdout, tpl[0], invoices, nil)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return globalRepository.Read()
		},
	}
	short = cmd.Flags().BoolP("short", "s", false, "show only names and IDs")

	return cmd
}
