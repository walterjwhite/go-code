package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/security/random"
)

var (
	lengthFlag       = flag.Int("l", 32, "Length of random string to produce")
	characterSetFlag = flag.String("c", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Character Set")
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	fmt.Println(random.StringWithCharset(*lengthFlag, *characterSetFlag))
}
