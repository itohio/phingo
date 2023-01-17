package app

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) initPreferences() {
}

func (a *App) ShowPreferences() {
	dlg := dialog.NewForm(
		"Settings",
		"OK", "Cancel",
		[]*widget.FormItem{},
		func(b bool) {
			if !b {
				return
			}

		},
		a.window,
	)

	dlg.Resize(a.window.Canvas().Size())
	dlg.Show()
}
