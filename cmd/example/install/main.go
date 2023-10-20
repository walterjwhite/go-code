package main

import (
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/os/install"
)

func init() {
	application.Configure()
}

func main() {
	install.SystemInstaller.Install("vim")
	install.NPMinstaller.Install("castnow")
}
