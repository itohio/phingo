package app

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	tabInvoicesName  = "Invoices"
	tabProjectsName  = "Projects"
	tabClientsName   = "Clients"
	tabAccountsName  = "Accounts"
	tabTemplatesName = "Templates"
)

func (a *App) makeMenu() *widget.Toolbar {

	addAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		switch a.tabs.Selected().Text {
		case tabInvoicesName:
			a.newInvoice()
		case tabProjectsName:
			a.newProject()
		case tabClientsName:
			a.newClient()
		case tabAccountsName:
			a.newAccount()
		case tabTemplatesName:
		}
	})

	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		a.ShowPreferences()
	})
	saveAction := widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		dialog.NewConfirm(
			"Writing database",
			"You are about to write to the database. This will overwrite the database with your local copy. Are you sure?",
			func(b bool) {
				if !b {
					return
				}
				a.repo.Write()
			},
			a.window,
		)
	})
	loadAction := widget.NewToolbarAction(theme.DownloadIcon(), func() {
		dialog.NewConfirm(
			"Loading database",
			"You are about to load the database. This will overwrite your local copy with the database. Are you sure?",
			func(b bool) {
				if !b {
					return
				}
				a.repo.Read()
			},
			a.window,
		)
	})

	actions := []widget.ToolbarItem{
		addAction,
		widget.NewToolbarSpacer(),
		saveAction,
		loadAction,
		settingsAction,
		widget.NewToolbarSeparator(),
	}

	return widget.NewToolbar(actions...)
}
