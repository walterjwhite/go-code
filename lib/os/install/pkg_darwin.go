package install

import (
	"context"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/runner"
	"time"
)

type AppleInstaller struct {
	root    string
	options string

	installTimeout   time.Duration
	bootstrapTimeout time.Duration
}

func (i *AppleInstaller) Install(packageNames ...string) {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *AppleInstaller) Uninstall(packageNames ...string) {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *AppleInstaller) IsInstalled(packageName string) bool {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *AppleInstaller) Bootstrap() {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *AppleInstaller) Cleanup() {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func (i *AppleInstaller) Update() {
	logging.Panic(fmt.Errorf("Not Implemented"))
}

func init() {
	SystemInstaller = &AppleInstaller{}
	SystemInstaller.Bootstrap()
}
