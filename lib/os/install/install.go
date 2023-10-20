package install

import (
	"context"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/os/sudo"
	"github.com/rs/zerolog/log"
	"time"
)

type Installer interface {
	Install(packageName ...string)
	Uninstall(packageName ...string)
	IsInstalled(packageName string) bool
	Bootstrap()
	Update()

	Cleanup()

	BootstrapType(typeName string)
}

var SystemInstaller Installer

func checkStatus(timeout time.Duration, cmd string, arguments ...string) bool {
	log.Debug().Msgf("timeout: %v", timeout)

	log.Debug().Msgf("cmd: %v, args: %v", cmd, arguments)

	ctx, cancel := context.WithTimeout(application.Context, timeout)
	defer cancel()

	status, err := sudo.Run(ctx, cmd, arguments...)
	logging.Panic(err)

	return status == 0
}
