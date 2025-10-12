package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/security/random"
)

var (
	lengthFlag       = flag.Int("l", 32, "Length of random string to produce")
	characterSetFlag = flag.String("chars", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Character Set")
)

func init() {
	application.Configure()
}

func main() {
	fmt.Println(random.StringWithCharset(*lengthFlag, *characterSetFlag))
}
