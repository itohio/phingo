package types

import (
	"fmt"
	"strings"

	"github.com/itohio/phingo/pkg/bi"
)

type InvoicesArr []*Invoice

func (arr InvoicesArr) ById(id string) *Invoice {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (arr InvoicesArr) ByCode(code string) *Invoice {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Code == code {
			return a
		}
	}
	return nil
}

func (arr InvoicesArr) ByYear(year int) *Invoice {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Year() == year {
			return a
		}
	}
	return nil
}

func (inv *Invoice) MakeId() string {
	return fmt.Sprintf("%d-%s", inv.Year(), strings.ReplaceAll(inv.Code, " ", "_"))
}

func (inv *Invoice) Year() int {
	t, err := bi.Parse(inv.IssueDate)
	if err != nil {
		return 0
	}
	return t.Year()
}

func (inv *Invoice) MakeFileName() string {
	format := "{Full Name}_{Code}_{Issue Date}"
	prj := inv.Project
	cl := inv.Client
	acc := inv.Account
	switch {
	case cl != nil && cl.InvoiceFileNameFormat != "":
		format = cl.InvoiceFileNameFormat
	case acc != nil && acc.InvoiceFileNameFormat != "":
		format = acc.InvoiceFileNameFormat
	}
	tokens := map[string]func() string{
		"{Full Name}": func() string {
			if acc == nil {
				return "-"
			}
			if val, ok := acc.Contact[ContactFullName]; ok {
				return val
			}
			return "-"
		},
		"{Name}": func() string {
			if acc == nil {
				return "-"
			}
			if val, ok := acc.Contact[ContactName]; ok {
				return val
			}
			return "-"
		},
		"{Account Name}": func() string {
			if acc == nil {
				return "-"
			}
			return acc.Name
		},
		"{Client Name}": func() string {
			if cl == nil {
				return "-"
			}
			return cl.Name
		},
		"{Project Name}": func() string {
			if prj == nil {
				return "-"
			}
			return prj.Name
		},
		"{Personal Code}": func() string {
			if acc == nil {
				return "-"
			}
			if val, ok := acc.Contact[ContactPersonalCode]; ok {
				return val
			}
			return "-"
		},
		"{Reg. #}": func() string {
			if acc == nil {
				return "-"
			}
			if val, ok := acc.Contact[ContactRegistrationNumber]; ok {
				return val
			}
			return "-"
		},
		"{Company Reg. #}": func() string {
			if acc == nil {
				return "-"
			}
			if val, ok := acc.Contact[ContactCompanyRegistrationNumber]; ok {
				return val
			}
			return "-"
		},
		"{Client Reg. #}": func() string {
			if cl == nil {
				return "-"
			}
			if val, ok := cl.Contact[ContactRegistrationNumber]; ok {
				return val
			}
			return "-"
		},
		"{Client Company Reg. #}": func() string {
			if cl == nil {
				return "-"
			}
			if val, ok := cl.Contact[ContactCompanyRegistrationNumber]; ok {
				return val
			}
			return "-"
		},
		"{Code}": func() string {
			return inv.Code
		},
		"{Issue Date}": func() string {
			t, err := bi.Parse(inv.IssueDate)
			if err != nil {
				return "-"
			}
			return t.Format("20061011")
		},
		"{Due Date}": func() string {
			t, err := bi.Parse(inv.DueDate)
			if err != nil {
				return "-"
			}
			return t.Format("20061011")
		},
	}
	for k, val := range tokens {
		format = strings.ReplaceAll(format, k, val())
	}
	return SanitizePath(format)
}
