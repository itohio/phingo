package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newClientsContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"Name", 150},
		{"Description", 150},
		{"Invoice File Format", 150},
		{"ID", 150},
	}
	ret := a.newTable(
		func() int { return len(a.clients) },
		headers,
		func(row, col int, l *widget.Label) string {
			cl := a.clients[row]
			switch col {
			case 0:
				return cl.Name
			case 1:
				return cl.Description
			case 2:
				return cl.InvoiceFileNameFormat
			case 3:
				return cl.Id
			}
			return "-"
		},
		nil,
	)

	return ret
}
