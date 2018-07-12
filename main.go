package main

import (
	"os"

	"github.com/gdamore/tcell"
	"github.com/jeffsaracco/goctobox/octobox"
	"github.com/rivo/tview"
)

var notifications *[]octobox.Notification

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	apiToken := os.Getenv("OCTOBOX_API_TOKEN")
	octoboxClient := octobox.New("https://octobox.githubapp.com", apiToken)

	fillTable(octoboxClient, table)
	table.SetSelectable(true, false)
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})
	table.SetSelectedFunc(handleSelection(table, notifications))
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			fillTable(octoboxClient, table)
		}
		return event
	})

	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}

}

func fillTable(octoboxClient *octobox.Client, table *tview.Table) {
	table.Clear()
	notifications = octoboxClient.GetNotifications()

	cols, rows := 4, len(*notifications)

	for r := 0; r < rows; r++ {
		setTitleRow(table)

		for c := 0; c < cols; c++ {
			notification := (*notifications)[r]
			setNotificationRow(table, notification, r+1, c)
		}
	}
}

func contains(s []int, e int) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}

func handleSelection(table *tview.Table, notifications *[]octobox.Notification) func(int, int) {
	selected := []int{}
	return func(row int, column int) {
		if row == 0 {
			return
		}
		notif := (*notifications)[row-1]
		cols := table.GetColumnCount()
		color := tcell.ColorRed

		if i, ok := contains(selected, notif.ID); ok {
			copy(selected[i:], selected[i+1:])
			selected[len(selected)-1] = 0
			selected = selected[:len(selected)-1]
			color = tcell.ColorWhite
		} else {
			selected = append(selected, notif.ID)
		}

		for c := 0; c < cols; c++ {
			table.GetCell(row, c).SetTextColor(color)
		}
	}
}

func getColumnData(notification octobox.Notification, column int) string {
	switch column {
	case 0:
		return notification.Subject.Type
	case 1:
		return notification.Subject.Title
	case 2:
		return notification.Repo.Name
	case 3:
		return notification.Reason
	}

	return ""
}

func setTitleRow(table *tview.Table) {
	columns := []string{"Type", "Title", "Repo", "Reason"}
	color := tcell.ColorBlue
	for index, col := range columns {
		table.SetCell(
			0,
			index,
			tview.NewTableCell(col).SetTextColor(color).SetAlign(tview.AlignCenter),
		)
	}
}

func setNotificationRow(table *tview.Table, notification octobox.Notification, row, column int) {
	color := tcell.ColorWhite
	if notification.Unread {
		color = tcell.ColorGreen
	}
	data := getColumnData(notification, column)
	table.SetCell(
		row,
		column,
		tview.NewTableCell(data).SetTextColor(color).SetAlign(tview.AlignCenter),
	)
}
