package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/random"
)

var (
	lengthFlag       = flag.Int("length", 32, "Length of random string to produce")
	characterSetFlag = flag.String("charset", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Character Set")
)

func init() {
	application.Configure()
}

func main() {
	fmt.Println(random.StringWithCharset(*lengthFlag, *characterSetFlag))
}
