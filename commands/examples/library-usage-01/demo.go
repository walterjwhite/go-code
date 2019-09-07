package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
)

func main() {
	ctx := application.Configure()
	application.Wait(ctx)
}
