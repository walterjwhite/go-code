package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	//"github.com/walterjwhite/go-hi"
	//gotest "github.com/walterjwhite/go-test"
	"github.com/walterjwhite/gotest"
)

func main() {
	application.Configure()

	gotest.Test()
}
