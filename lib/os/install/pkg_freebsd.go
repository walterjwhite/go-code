package install

import (
	"time"
)

type FreeBSDPKGInstaller struct {
	root    string
	options string

	installTimeout   time.Duration
	bootstrapTimeout time.Duration
}


func (i *FreeBSDPKGInstaller) Install(packageNames ...string) {
	var args []string
	args = append(args, i.options)
	args = append(args, "install")
	args = append(args, "-y")
	args = append(args, packageNames...)

	_ = checkStatus(i.installTimeout, "pkg", args...)
}

func (i *FreeBSDPKGInstaller) Uninstall(packageNames ...string) {
	var args []string
	args = append(args, i.options)
	args = append(args, "delete")
	args = append(args, "-y")
	args = append(args, packageNames...)

	_ = checkStatus(i.installTimeout, "pkg", args...)
}

func (i *FreeBSDPKGInstaller) IsInstalled(packageName string) bool {
	return checkStatus(i.installTimeout, "pkg", i.options, "info", "-e", packageName)
}

func (i *FreeBSDPKGInstaller) Bootstrap() {
	if len(i.root) > 0 {
		_ = checkStatus(i.bootstrapTimeout, "mount", "/var/cache/pkg", i.root+"/var/cache/pkg")
	}
}

func (i *FreeBSDPKGInstaller) Cleanup() {
	if len(i.root) > 0 {
		_ = checkStatus(i.bootstrapTimeout, "umount", i.root+"/var/cache/pkg")
	}
}

func (i *FreeBSDPKGInstaller) Update() {
	_ = checkStatus(i.installTimeout, "pkg", "update", i.options)
}

func (i *FreeBSDPKGInstaller) BootstrapType(typeName string) {
	i.Install(typeName)
}

func init() {
	SystemInstaller = &FreeBSDPKGInstaller{installTimeout: 1 * time.Minute, bootstrapTimeout: 5 * time.Minute}
	SystemInstaller.Bootstrap()
}
