package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

func main() {
	application.Configure()
	//application.Wait(ctx)

	screenshot.Take("label", "detail")
}
