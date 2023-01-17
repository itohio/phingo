package app

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/repository"
)

type App struct {
	sync.Mutex
	app    fyne.App
	window fyne.Window
	repo   repository.Repository
}

var (
	appFQDN = "itohio.phingo"
)

func New(repo repository.Repository) *App {
	a := app.NewWithID(appFQDN)
	w := a.NewWindow("name")
	w.Resize(fyne.NewSize(500, 600))

	ret := &App{
		app:    a,
		window: w,
		repo:   repo,
	}

	ret.initPreferences()

	return ret
}

func (a *App) Run() error {
	a.makeContent()
	a.window.ShowAndRun()
	return nil
}

func (a *App) makeContent() {
	a.window.SetContent(
		container.NewBorder(
			a.makeMenu(), nil, // Top, Bottom
			nil, nil, // Left, Right
			container.NewAppTabs(
				container.NewTabItem("Invoices", widget.NewAccordion()),
				container.NewTabItem("Clients", widget.NewAccordion()),
				container.NewTabItem("Accounts", widget.NewAccordion()),
			),
		),
	)
}
