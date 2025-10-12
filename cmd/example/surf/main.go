package main

import (
	"fmt"
	"gopkg.in/headzoo/surf.v1"
)

func main() {
	bow := surf.NewBrowser()
	err := bow.Open("http://golang.org")
	if err != nil {
		panic(err)
	}

	fmt.Println(bow.Title())
}
