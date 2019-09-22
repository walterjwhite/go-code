package service

import (
	"../data"
	"../health"
	"strings"
)

var header = data.Header{[]string{"Service", "Listening Addresses"}}

func Convert(services []Service) []data.Row {
	rows := make([]data.Row, 0)

	for i := 0; i < len(services); i++ {
		row := data.Row{getStatus(services[i].Health), getRow(&services[i])}
		rows = append(rows, row)
	}

	return rows
}

func getStatus(status int) int {
	if status == health.HEALTH_GOOD {
		return data.Good
	}
	if status < health.HEALTH_GOOD {
		return data.Warning
	}

	return data.Bad
}

func getRow(service *Service) []string {
	return []string{service.Name /*service.Status, */, strings.Join(service.Addresses, ",")}
}
