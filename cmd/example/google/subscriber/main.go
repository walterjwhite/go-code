package main

import (
	"os"
	"github.com/walterjwhite/go-code/lib/net/google"
)

func main() {
	google.Subscribe(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
}
