package install

import (
      "os/exec"
)

func IsCommandAvailable(name string) bool {
      cmd := exec.Command("command", "-v", name)
      err := cmd.Run()
	  return err == nil
}
