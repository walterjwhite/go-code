package pf

import (
	"fmt"
	"github.com/walterjwhite/go/lib/application/logging"
	"os/exec"
)

var (
	tableName string
)

func pfAdd(ip string) {
	_pfAction(ip, "add")
}

func pfRemove(ip string) {
	_pfAction(ip, "delete")
}

func _pfAction(ip, action string) {
	cmd := exec.Command("pfctl", "-t", tableName, "-T", action, ip)
	logging.Warn(cmd.Run(), false, fmt.Sprintf("Error %v %v", action, ip))
}
