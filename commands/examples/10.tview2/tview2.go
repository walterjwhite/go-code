package main

import (
	"github.com/rivo/tview"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var a tview.Application

func onChange( /*a *tview.Application*/ ) {
	a.Draw()
}

func init() {
	application.Configure()
}

func main() {
	newPrimitive := func(text string) *tview.TextView {
		return tview.NewTextView().
			//SetTextAlign(tview.AlignCenter).
			SetTextAlign(tview.AlignLeft).
			SetText(text).
			SetChangedFunc(onChange)
	}
	services := newPrimitive("Services\natd: UP\nargus: UP\nchronyd: UP 192.168.55.1:123 192.168.56.1:123\nunbound: UP 192.168.55.1:53 192.168.56.1:53\nsshd: UP 192.168.55.1:22 192.168.56.1:22\nnetif.lw: UP\nnetif.ll: UP\nnetif.wl: UP\ndhcpd: UP\ndnscrypt-proxy: UP\ndnstap: UP\nhostapd: UP\nprivoxy: UP 192.168.55.1:8118 192.168.56.1:8118\nshorewall: UP\nsquid: UP 127.0.0.1:3128\ntor: UP 127.0.0.1:9050\nunbound: UP 192.168.55.1:53")

	interfaces := newPrimitive("lw: 173.75.130.27\nll: 192.168.55.1\nwl: 192.168.56.1\n")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(70, 30).
		//SetBorders(true).
		AddItem(newPrimitive("2019/03/22 10:45 AM   built on: 2019/03/01 11:01 PM  up: 4 days"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("syslog / dmesg output"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(services, 1, 0, 1, 1, 0, 100, false).
		AddItem(interfaces, 1, 1, 1, 1, 0, 100, false)

	// update the text in parallel (working)
	/*
		go func() {
			index := 0
			for index < 1000 {
				//main.SetText(main.GetText(true) + fmt.Sprintf("(%d)", index))
				//main.SetText(fmt.Sprintf("(%d)", index))
				main.Clear().Write([]byte(fmt.Sprintf("(%d)", index)))
				index += 1

				time.Sleep(1 * time.Second)
			}
		}()
	*/

	logging.Panic(a.SetRoot(grid, true).Run())
}

func init() {
	a = *tview.NewApplication()
}
