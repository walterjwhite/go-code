package install

import (
	"github.com/walterjwhite/go-code/lib/application/logging"

	"fmt"

	"time"
)

type NPMInstaller struct {
	options string

	installTimeout time.Duration
	isBootstrapped bool
}

var NPMinstaller *NPMInstaller

func (i *NPMInstaller) Install(packageName string) {
	NPMinstaller.Bootstrap()
	_ = checkStatus(i.installTimeout, "npm", i.options, "install", "-g", packageName)
}

func (i *NPMInstaller) Uninstall(packageName string) {
	NPMinstaller.Bootstrap()
	_ = checkStatus(i.installTimeout, "npm", i.options, "uninstall", "-g", packageName)
}

func (i *NPMInstaller) IsInstalled(packageName string) bool {
	NPMinstaller.Bootstrap()

	return checkStatus(i.installTimeout, "npm", i.options, "list", "-g", packageName)
}

func (i *NPMInstaller) Bootstrap() {
	if i.isBootstrapped {
		return
	}

	SystemInstaller.Bootstrap()

	SystemInstaller.Install("npm")
	NPMinstaller.Bootstrap()
	i.isBootstrapped = IsCommandAvailable("npm")

	if !i.isBootstrapped {
		logging.Panic(fmt.Errorf("error bootstrapping NPM"))
	}
}

func (i *NPMInstaller) Cleanup() {
	logging.Panic(fmt.Errorf("not implemented"))
}

func (i *NPMInstaller) Update() {
	logging.Panic(fmt.Errorf("not implemented"))
}

func init() {
	NPMinstaller = &NPMInstaller{installTimeout: 1 * time.Minute, isBootstrapped: IsCommandAvailable("npm")}
}
