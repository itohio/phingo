package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newProjectsContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"Name", 150},
		{"Description", 150},
		{"Start", 150},
		{"Deadline", 150},
		{"ID", 150},
	}
	ret := a.newTable(
		func() int { return len(a.projects) },
		headers,
		func(row, col int, l *widget.Label) string {
			prj := a.projects[row]
			switch col {
			case 0:
				return prj.Name
			case 1:
				return prj.Description
			case 2:
				return prj.StartDate
			case 3:
				return prj.EndDate
			case 4:
				return prj.Id
			}
			return "-"
		},
		nil,
	)

	return ret
}
