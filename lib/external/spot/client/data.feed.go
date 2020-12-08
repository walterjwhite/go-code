package client

import (
	"net/http"
)

type Feed struct {
	Id                   string
	Name                 string
	Description          string
	Status               string
	Usage                int
	DaysRange            int
	DetailedMessageShown bool
	Type                 string

	client *http.Client
}
