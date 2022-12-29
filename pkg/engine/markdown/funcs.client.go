package engine

import (
	"fmt"
	"text/template"

	"github.com/itohio/phingo/pkg/types"
)

func makeNotesFunc(config *types.Config) func(*types.Client) []string {
	return func(val *types.Client) []string {
		return val.Notes
	}
}

func makeClientFunc(config *types.Config) func(*types.Client) []string {
	return func(val *types.Client) []string {
		s := []string{
			fmt.Sprintf("**%s**", config.Locale.Translate("CLIENT")),
		}
		s = append(s, makeContactsFunc(config)(val.Contact)...)

		return s
	}
}

func makeClientFuncs(context *types.ClientTemplateContext) template.FuncMap {
	return template.FuncMap{
		"Contacts": makeContactsFunc(context.Config),
		"Notes":    makeNotesFunc(context.Config),
		"Client":   makeClientFunc(context.Config),
	}
}
