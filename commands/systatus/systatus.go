package main

import (
	"./data"
	"./disk"
	"./interfaces"
	"./service"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"time"

	//"./health"
	//"./tail"
	"./irpc"
	"./ui"
	"./ui/table"
	"./ui/text"
)

var clientInstance = irpc.New()

type RefreshableTable struct {
	Table   *tview.Table
	Header  *data.Header
	Refresh func() []data.Row
}

type RefreshableText struct {
	Text    *tview.TextView
	Label   string
	Refresh func() string
}

var refreshableTables []RefreshableTable
var refreshableTexts []RefreshableText

//func createText(

const refreshInterval = 1 * time.Minute

var application tview.Application

func refreshAll() error {
	refreshTables()
	refreshTexts()

	application.Draw()

	return nil
}

// TODO: execute these in parallel
func refreshTables() {
	for i := 0; i < len(refreshableTables); i++ {
		table.Refresh(refreshableTables[i].Table, refreshableTables[i].Refresh())
	}
}

// TODO: execute these in parallel
func refreshTexts() {
	for i := 0; i < len(refreshableTexts); i++ {
		text.Refresh(refreshableTexts[i].Text, refreshableTexts[i].Label, refreshableTexts[i].Refresh())
	}
}

func periodic(fn func() error) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	fn()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				log.Fatalf("Error executing Periodic %v", err)
			}
		}
	}
}

func buildHeader(grid *tview.Grid) {
	grid.AddItem(refreshableTexts[1].Text, 0, 0, 1, 1, 2, 0, false)
	grid.AddItem(refreshableTexts[2].Text, 0, 1, 1, 1, 2, 0, false)
	grid.AddItem(refreshableTexts[3].Text, 0, 2, 1, 1, 2, 0, false)
}

func buildMiddle(grid *tview.Grid) {
	grid.AddItem(refreshableTables[0].Table, 1, 0, 3, 2, 20, 100, false)

	grid.AddItem(refreshableTables[1].Table, 1, 2, 1, 1, 0, 50, false)
	grid.AddItem(refreshableTables[2].Table, 2, 2, 1, 1, 0, 50, false)
	//grid.AddItem(refreshableTables[3].Table, 3, 2, 1, 1, 0, 50, false)
}

func buildBottom(grid *tview.Grid) {
	grid.AddItem(refreshableTexts[0].Text, 4, 0, 1, 3, 5, 0, false)
}

func main() {
	grid := tview.NewGrid().SetRows(2, 0, 0)

	buildHeader(grid)
	buildMiddle(grid)
	buildBottom(grid)

	//refresh underlying data
	go periodic(refreshAll)

	if err := application.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

func ignoreEvents(event *tcell.EventKey) *tcell.EventKey {
	// do NOTHING
	return nil
}

func init() {
	application = *tview.NewApplication()
	//application.SetInputCapture(ignoreEvents)

	refreshableTables = make([]RefreshableTable, 0)
	refreshableTables = append(refreshableTables, RefreshableTable{table.Table(), service.Header(), clientInstance.Services})
	refreshableTables = append(refreshableTables, RefreshableTable{table.Table(), disk.Header(), clientInstance.Disks})
	refreshableTables = append(refreshableTables, RefreshableTable{table.Table(), interfaces.Header(), clientInstance.Interfaces})

	//refreshableTables = append(refreshableTables, RefreshableTable{table.Table(), network.Header(), network.Data})

	// initialize headers
	for i := 0; i < len(refreshableTables); i++ {
		table.InitHeader(refreshableTables[i].Table, refreshableTables[i].Header)
	}

	refreshableTexts = make([]RefreshableText, 0)
	// syslog
	refreshableTexts = append(refreshableTexts, RefreshableText{text.Text(), "syslog", clientInstance.Logs})

	refreshableTexts = append(refreshableTexts, RefreshableText{text.Text(), "Date/Time", ui.CurrentDateTime})
	refreshableTexts = append(refreshableTexts, RefreshableText{text.Text(), "Build Date/Time", clientInstance.BuildDateTime})
	refreshableTexts = append(refreshableTexts, RefreshableText{text.Text(), "Uptime", clientInstance.Uptime})
}
