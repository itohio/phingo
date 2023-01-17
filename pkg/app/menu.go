package app

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (a *App) makeMenu() *widget.Toolbar {

	settingsAction := widget.NewToolbarAction(theme.SettingsIcon(), func() {
		a.ShowPreferences()
	})

	var actions []widget.ToolbarItem

	actions = append(
		actions,
		[]widget.ToolbarItem{
			widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			}),
			widget.NewToolbarAction(theme.DownloadIcon(), func() {
			}),
		}...,
	)

	actions = append(
		actions,
		[]widget.ToolbarItem{
			widget.NewToolbarSpacer(),
			settingsAction,
			widget.NewToolbarSeparator(),
		}...,
	)

	return widget.NewToolbar(actions...)
}
