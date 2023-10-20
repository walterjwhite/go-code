package install

import (


	"github.com/walterjwhite/go-code/lib/application/logging"

	"fmt"

	"time"
)

type GoInstaller struct {
	// root    string
	options string

	installTimeout   time.Duration
	// bootstrapTimeout time.Duration
}

func (i *GoInstaller) Install(packageName string) {
	_ = checkStatus(i.installTimeout, "go", "install", i.options, packageName)
}

func (i *GoInstaller) Uninstall(packageName string) {
	_ = checkStatus(i.installTimeout, "go", "uninstall", packageName)
}

func (i *GoInstaller) IsInstalled(packageName string) bool {
	logging.Panic(fmt.Errorf("Not Implemented"))

	return false
}

func (i *GoInstaller) Bootstrap() {
	SystemInstaller.Bootstrap()

	SystemInstaller.Install("go")
}

func (i *GoInstaller) Cleanup() {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *GoInstaller) Update() {
	logging.Panic(fmt.Errorf("Not Implemented"))
}
