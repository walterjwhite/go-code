package health

import (
	"os/exec"
)

func IsShorewallRunning() int {
	cmd := exec.Command("/usr/sbin/shorewall", "status")
	_, err := cmd.CombinedOutput()

	if err == nil {
		return HEALTH_GOOD
	}

	return HEALTH_BAD
}
