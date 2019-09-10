package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

func main() {
	application.Configure()
	defer application.OnCompletion()

	//application.Wait(ctx)

	path.WithSessionDirectory("~/.audit")
	screenshot.Take("label", "detail")
}
