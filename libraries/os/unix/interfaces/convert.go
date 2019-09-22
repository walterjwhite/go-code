package interfaces

import (
	"github.com/walterjwhite/go-application/libraries/supervisor/data"
)

const (
	UP   = "Up"
	DOWN = "Down"
)

var header = data.Header{[]string{"Interface", "Status", "IP"}}

func Convert(interfaces []Interface) []data.Row {
	rows := make([]data.Row, 0)

	for i := 0; i < len(interfaces); i++ {
		rows = append(rows, data.Row{getStatus(interfaces[i].Up), getRow(interfaces[i])})
	}

	return rows
}

func getStatus(isUp bool) int {
	if isUp {
		return data.Good
	}

	return data.Bad
}

func getRow(networkInterface Interface) []string {
	return []string{networkInterface.Name, getStatusString(networkInterface.Up), networkInterface.IP}
}

func getStatusString(up bool) string {
	if up {
		return UP
	}

	return DOWN
}
