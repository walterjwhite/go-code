package main

import (
	"github.com/walterjwhite/go-code/lib/net/google"
	"os"
)

func main() {
	google.Publish(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
}
