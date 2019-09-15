package main

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
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
				if strings.Index(lorem[word], "started") >= 0 {
					tableCell.SetBackgroundColor(tcell.ColorGreen)
				} else if strings.Index(lorem[word], "stopping") >= 0 || strings.Index(lorem[word], "starting") >= 0{
					tableCell.SetBackgroundColor(tcell.ColorYellow)
				} else {
					tableCell.SetBackgroundColor(tcell.ColorRed)
				}
			}
			
			table.SetCell(r, c, tableCell)
			
			word = (word + 1) % len(lorem)
		}
	}
	
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
