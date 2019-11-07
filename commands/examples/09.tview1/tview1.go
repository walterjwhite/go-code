package main

import (
	"github.com/rivo/tview"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func init() {
	application.Configure()
}

func main() {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	logging.Panic(tview.NewApplication().SetRoot(box, true).Run())
}
