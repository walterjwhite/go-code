package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/secrets"
)

// TODO: add support for flags
// instead of specifying the key type (email, user, pass), use a flag instead (-e, -u, -p)
func init() {
	application.Configure()
}

func main() {
	fmt.Println(secrets.Base64(secrets.DoEncrypt([]byte(flag.Args()[0]))))
}
