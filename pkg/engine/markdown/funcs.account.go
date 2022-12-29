package engine

import (
	"fmt"
	"text/template"

	"github.com/itohio/phingo/pkg/types"
)

func makeContactsFunc(config *types.Config) func(map[string]string) []string {
	return func(val map[string]string) []string {
		res := []string{}
		keys := []string{
			"Full Name",
			"Name",
			"Personal Code",
			"Reg. #",
			"Company Reg. #",
			"VAT",
			"Address",
			"Phone",
			"Cell",
			"Email",
			"Bank Account",
			"Wallet Address",
			"Director",
		}
		if val == nil {
			return res
		}
		name := false
		for i, key := range keys {
			if name && i < 2 {
				continue
			}
			if v, ok := val[key]; ok {
				contact := v
				if key != keys[0] && key != keys[1] {
					contact = fmt.Sprintf("%s: %s", config.Locale.Translate(key), v)
				} else {
					name = true
				}
				res = append(res, contact)
			}
		}

		return res
	}
}

func makeAccountFunc(config *types.Config) func(*types.Account) []string {
	return func(val *types.Account) []string {
		s := []string{
			fmt.Sprintf("**%s**", config.Locale.Translate("SELLER")),
		}
		s = append(s, makeContactsFunc(config)(val.Contact)...)
		return s
	}
}

func makeAccountFuncs(context *types.AccountTemplateContext) template.FuncMap {
	return template.FuncMap{
		"Contacts": makeContactsFunc(context.Config),
		"Account":  makeAccountFunc(context.Config),
	}
}
