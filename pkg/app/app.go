package app

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/repository"
	"github.com/itohio/phingo/pkg/types"
)

type App struct {
	sync.Mutex
	app    fyne.App
	window fyne.Window
	tabs   *container.AppTabs
	repo   repository.Repository

	accounts  []*types.Account
	clients   []*types.Client
	projects  []*types.Project
	invoices  []*types.Invoice
	templates []*types.Template

	filtersInvoices []string
}

var (
	appFQDN = "itohio.phingo"
)

func New(repo repository.Repository) *App {
	a := app.NewWithID(appFQDN)
	w := a.NewWindow("name")
	w.Resize(fyne.NewSize(800, 600))

	ret := &App{
		app:    a,
		window: w,
		repo:   repo,
	}

	ret.readRepo()
	ret.initPreferences()

	return ret
}

func (a *App) readRepo() {
	a.accounts = a.repo.Accounts()
	a.clients = a.repo.Clients()
	a.projects = a.repo.Projects()
	a.templates = a.repo.Templates()
	a.invoices = a.repo.Invoices(a.filtersInvoices...)
}

func (a *App) Run() error {
	a.makeContent()
	a.window.ShowAndRun()
	return nil
}

func (a *App) makeContent() {
	a.tabs = container.NewAppTabs(
		container.NewTabItem(tabInvoicesName, a.newInvoicesContents()),
		container.NewTabItem(tabProjectsName, a.newProjectsContents()),
		container.NewTabItem(tabClientsName, a.newClientsContents()),
		container.NewTabItem(tabAccountsName, a.newAccountsContents()),
		container.NewTabItem(tabTemplatesName, a.newTemplatesContents()),
	)

	a.window.SetContent(
		container.NewBorder(
			a.makeMenu(), nil, // Top, Bottom
			nil, nil, // Left, Right
			a.tabs,
		),
	)
}

func (a *App) setColumnWidths(table *widget.Table, headers []tableHeader) {
	for i := range headers {
		table.SetColumnWidth(i, headers[i].width)
	}
}

type tableHeader struct {
	name  string
	width float32
}

func (a *App) newTable(rows func() int, headers []tableHeader, upd func(row, col int, co *widget.Label) string, onSel func(row int)) *widget.Table {
	ret := widget.NewTable(
		func() (int, int) {
			return rows(), len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(".")
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			l := co.(*widget.Label)
			if tci.Row == 0 {
				l.TextStyle.Bold = true
				l.SetText(headers[tci.Col].name)
				return
			}
			l.TextStyle.Bold = false
			l.SetText(upd(tci.Row-1, tci.Col, l))
		},
	)
	a.setColumnWidths(ret, headers)

	if onSel != nil {
		ret.OnSelected = func(id widget.TableCellID) {
			if id.Row == 0 {
				return
			}
			onSel(id.Row - 1)
		}
	}

	return ret
}
