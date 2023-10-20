package secrets

import (
	"os/exec"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Get(secretName string) string {
	out, err := exec.Command("secrets", "get", secretName).Output()
    if len(out) == 0 {
		logging.Panic(err)
	}

    return string(out[:])
}