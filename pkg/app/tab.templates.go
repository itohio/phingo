package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newTemplatesContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"ID", 150},
		{"What", 150},
	}
	ret := a.newTable(
		func() int { return len(a.templates) },
		headers,
		func(row, col int, l *widget.Label) string {
			prj := a.templates[row]
			switch col {
			case 0:
				return prj.Id
			case 1:
				return prj.What
			}
			return "-"
		},
		nil,
	)

	return ret
}
