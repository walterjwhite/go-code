package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"strings"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var a tview.Application

func onChange( /*a *tview.Application*/ ) {
	a.Draw()
}

func buildServiceTable() *tview.Table {
	table := tview.NewTable()
	content := strings.Split("Service Status IP:Port unbound started 192.168.55.1:53,192.168.56.1:53 atd stopped N/A dhcpd crashed N/A dnscrypt-proxy started 127.0.0.1:5300 haveged started N/A netif.lw started N/A sshd started 127.0.0.1:22,192.168.55.2:22,192.168.55.1:22,192.168.56.1:22 cupsd starting N/A netif.ll started N/A netif.wl started N/A privoxy.tor started 192.168.55.1:8119,192.168.56.1:8119 argus started N/A chronyd started 192.168.55.1:123,192.168.56.1:123 dnstap started N/A hostapd started N/A privoxy starting 192.168.55.1:8118,192.168.56.1:8118 shorewall started N/A squid started 127.0.0.1:3128 tor started 127.0.0.1:9050", " ")
	cols, rows := 3, 20
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			tableCell := tview.NewTableCell(content[word]).
				SetAlign(tview.AlignLeft)

			if r == 0 {
				// SetAttributes sets the cell's text attributes. You can combine different
				// attributes using bitmask operations:
				//
				//   cell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
				tableCell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
			}

			if c == 1 && r > 0 {
				if strings.Contains(content[word], "started") {
					tableCell.SetBackgroundColor(tcell.ColorGreen)
				} else if strings.Contains(content[word], "stopping") || strings.Contains(content[word], "starting") {
					tableCell.SetBackgroundColor(tcell.ColorYellow)
				} else {
					tableCell.SetBackgroundColor(tcell.ColorRed)
				}

				tableCell.SetTextColor(tcell.ColorBlack)
			}

			table.SetCell(r, c, tableCell)

			word = (word + 1) % len(content)
		}
	}

	return table
}

func buildInterfaceTable() *tview.Table {
	table := tview.NewTable()
	content := strings.Split("Interface IP UP lw 1.1.1.1 UP ll 192.168.55.1 UP wl 192.168.56.1 DOWN", " ")
	cols, rows := 3, 4
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			tableCell := tview.NewTableCell(content[word]).
				SetAlign(tview.AlignLeft)

			if r == 0 {
				tableCell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
			}

			if c == 2 && r > 0 {
				if strings.Contains(content[word], "UP") {
					tableCell.SetBackgroundColor(tcell.ColorGreen)
				} else {
					tableCell.SetBackgroundColor(tcell.ColorRed)
				}

				tableCell.SetTextColor(tcell.ColorBlack)
			}

			table.SetCell(r, c, tableCell)
			word = (word + 1) % len(content)
		}
	}

	return table
}

func buildNetworkStatusTable() *tview.Table {
	table := tview.NewTable()
	content := strings.Split("Network Tool,Parameters,Stats,dig,google.com,8.8.8.8 (20ms),ping,8.8.8.8,200ms (100%),wget,google.com,1s (OK)", ",")
	cols, rows := 3, 4
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			tableCell := tview.NewTableCell(content[word]).
				SetAlign(tview.AlignLeft)

			if r == 0 {
				tableCell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
			}

			table.SetCell(r, c, tableCell)
			word = (word + 1) % len(content)
		}
	}

	return table
}

func buildDiskTable() *tview.Table {
	table := tview.NewTable()
	content := strings.Split("Mount Point,%Usage,MB Free,/rw,25,128M,/var/cache/squid,50,4096M,/var/log,3,19G,/boot,82,42M", ",")
	cols, rows := 3, 5
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			tableCell := tview.NewTableCell(content[word]).
				SetAlign(tview.AlignLeft)

			if r == 0 {
				tableCell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
			}

			if c == 1 && r > 0 {
				i, err := strconv.ParseInt(content[word], 10, 16)
				if err != nil {
					panic("Error")
				}

				if i <= 25 {
					tableCell.SetBackgroundColor(tcell.ColorGreen)
				} else if i <= 50 {
					tableCell.SetBackgroundColor(tcell.ColorYellow)
				} else {
					tableCell.SetBackgroundColor(tcell.ColorRed)
				}

				tableCell.SetTextColor(tcell.ColorBlack)
			}

			table.SetCell(r, c, tableCell)
			word = (word + 1) % len(content)
		}
	}

	return table
}

func main() {
	application.Configure()

	newPrimitive := func(text string) *tview.TextView {
		return tview.NewTextView().
			//SetTextAlign(tview.AlignCenter).
			SetTextAlign(tview.AlignLeft).
			SetText(text).
			SetChangedFunc(onChange)
	}

	//interfaces := newPrimitive("lw: 173.75.130.27\nll: 192.168.55.1\nwl: 192.168.56.1\n")

	grid := tview.NewGrid().
		SetRows(2, 0, 0).
		//SetColumns(100, 0).
		//SetBorders(true).
		//(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool) *Grid {
		//AddItem(buildHeaderTable(), 0, 0, 1, 2, 2, 0, false).
		AddItem(newPrimitive("Date: 2019/03/22 10:45 AM"), 0, 0, 1, 1, 2, 0, false).
		AddItem(newPrimitive("Built on: 2019/03/01 11:01 PM"), 0, 1, 1, 1, 2, 0, false).
		AddItem(newPrimitive("up: 4 days"), 0, 2, 1, 1, 2, 0, false).
		AddItem(buildServiceTable(), 1, 0, 3, 2, 20, 100, false).
		AddItem(buildInterfaceTable(), 1, 2, 1, 1, 0, 50, false).
		AddItem(buildNetworkStatusTable(), 2, 2, 1, 1, 0, 50, false).
		AddItem(buildDiskTable(), 3, 2, 1, 1, 0, 50, false).
		AddItem(newPrimitive("syslog / dmesg output\nline1\nline2\nline3\nline4\nline5"), 4, 0, 1, 3, 5, 0, false)

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
