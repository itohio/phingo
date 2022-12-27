package repository

import (
	_ "embed"

	"github.com/itohio/phingo/pkg/types"
)

var (
	//go:embed projects.md
	defaultProjectsTemplate string

	//go:embed project.md
	defaultProjectTemplate string

	//go:embed invoices.md
	defaultInvoicesTemplate string

	//go:embed invoice.md
	defaultInvoiceTemplate string

	//go:embed accounts.md
	defaultAccountsTemplate string

	//go:embed clients.md
	defaultClientsTemplate string
)

func DefaultConfig() *types.Config {
	return &types.Config{
		Export: []*types.Config_Export{
			{
				What:   "console",
				Styles: []*types.Config_Style{},
			},
			{
				What:   "pdf",
				Styles: []*types.Config_Style{},
			},
		},
	}
}

func DefaultAccounts() *types.Accounts {
	return &types.Accounts{
		Accounts: []*types.Account{
			&types.Account{
				Id:    "default",
				Name:  "Default",
				Denom: "Eth",
				Contact: map[string]string{
					"Name":        "John Snow",
					"Email":       "John.Snow@email.com",
					"Address":     "The wall",
					"Eth address": "0xE774767569385C24740Faeffc0Ed6E1b4A87D619",
					"Denom":       "Eth",
					"Denom1":      "Eur",
				},
			},
		},
	}
}

func DefaultClients() *types.Clients {
	return &types.Clients{
		Clients: []*types.Client{
			&types.Client{
				Id:          "default",
				Name:        "Default",
				Description: "Default client",
				Contact: map[string]string{
					"Name":    "White Walker",
					"Email":   "white.walker@email.com",
					"Address": "Beyound The wall",
					"Denom":   "Ice",
				},
			},
		},
	}
}

func DefaultTemplates() []*types.Template {
	return []*types.Template{
		{
			FileName: "projects.md",
			Id:       "projects",
			What:     "projects",
			Text:     []byte(defaultProjectsTemplate),
		},
		{
			FileName: "project.md",
			Id:       "project",
			What:     "project",
			Text:     []byte(defaultProjectTemplate),
		},
		{
			FileName: "invoices.md",
			Id:       "invoices",
			What:     "invoices",
			Text:     []byte(defaultInvoicesTemplate),
		},
		{
			FileName: "invoice.md",
			Id:       "invoice",
			What:     "invoice",
			Text:     []byte(defaultInvoiceTemplate),
		},
		{
			FileName: "accounts.md",
			Id:       "accounts",
			What:     "accounts",
			Text:     []byte(defaultAccountsTemplate),
		},
		{
			FileName: "clients.md",
			Id:       "clients",
			What:     "clients",
			Text:     []byte(defaultClientsTemplate),
		},
	}
}

func DefaultProjects() []*types.Project {
	return nil
}

func DefaultInvoices() []*types.Invoice {
	return nil
}
