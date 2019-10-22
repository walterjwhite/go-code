package main

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func main() {
	application.Configure()

	app := tview.NewApplication()
	table := tview.NewTable()
	lorem := strings.Split("Service Status IP:Port unbound started 192.168.55.1:53,192.168.56.1:53 atd stopped N/A dhcpd crashed N/A dnscrypt-proxy started 127.0.0.1:5300", " ")
	cols, rows := 3, 5
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			tableCell := tview.NewTableCell(lorem[word]).
				SetAlign(tview.AlignLeft)

			if c == 1 && r > 0 {
				if strings.Contains(lorem[word], "started") {
					tableCell.SetBackgroundColor(tcell.ColorGreen)
				} else if strings.Contains(lorem[word], "stopping") || strings.Contains(lorem[word], "starting") {
					tableCell.SetBackgroundColor(tcell.ColorYellow)
				} else {
					tableCell.SetBackgroundColor(tcell.ColorRed)
				}
			}

			table.SetCell(r, c, tableCell)

			word = (word + 1) % len(lorem)
		}
	}

	logging.Panic(app.SetRoot(table, true).Run())
}
