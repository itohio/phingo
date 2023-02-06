package app

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/types"
)

func (a *App) makeDefaultContacts() map[string]string {
	ret := make(map[string]string, 10)
	for _, t := range []string{
		types.ContactName,
		types.ContactFullName,
		types.ContactAddress,
		types.ContactEmail,
		types.ContactBankAccount,
		types.ContactWalletAddress,
		types.ContactPhone,
	} {
		ret[t] = ""
	}
	return ret
}

func addContactItem(form *widget.Form, k, v string) {
	idx := len(form.Items)
	btn := widget.NewButtonWithIcon("",
		theme.DeleteIcon(),
		func() {
			form.Items = append(form.Items[:idx], form.Items[idx+1:]...)
			form.Refresh()
		},
	)
	eValue := widget.NewEntry()
	eValue.Text = v
	fi := widget.NewFormItem(
		k,
		container.NewBorder(nil, nil, nil, btn, eValue),
	)
	form.Items = append(form.Items, fi)
}

func (a *App) newContactItemAdder(form *widget.Form, kv map[string]string) func() {
	return func() {
		eType := widget.NewSelectEntry(types.DefaultContactTypes())
		eValue := widget.NewEntry()
		eType.Validator = validateChain(validateEmpty, validateUniqueKey(kv))
		dlg := dialog.NewForm(
			"Contact", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Type", eType),
				widget.NewFormItem("Value", eValue),
			},
			func(b bool) {
				if !b {
					return
				}
				kv[eType.Text] = eValue.Text
				addContactItem(form, eType.Text, eValue.Text)
				form.Refresh()
			},
			a.window,
		)
		es := dlg.MinSize()
		es.Width = 350
		dlg.Resize(es)
		dlg.Show()
	}
}
