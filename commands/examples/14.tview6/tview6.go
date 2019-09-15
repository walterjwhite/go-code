package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Left"), 0, 1, false).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (5 rows)"), 5, 1, false)).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
