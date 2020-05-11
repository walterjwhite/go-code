package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
)

func init() {
	application.Configure()
}

func main() {
	application.Wait()
}
