package pf

import (
	"os/exec"
)

var (
	tableName string
)

// func updatePFTable(tableName string){
// 	pfctl -t tableName -T replace -f <tempFile>
// }

func pfAdd(ip string) {
	_pfAction(ip, "add")
}

func pfRemove(ip string) {
	_pfAction(ip, "delete")
}

func _pfAction(ip, action string) {
	cmd := exec.Cmd("pfctl", "-t", tableName, "-T", action, ip)
	logging.Warn(cmd.Run())
}