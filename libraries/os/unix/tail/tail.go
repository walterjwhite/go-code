package tail

import (
	"log"
	"os/exec"
)

const (
	LINES_TO_READ = "-10"
	LOG_FILE      = "/var/log/messages"
)

func Data() string {
	cmd := exec.Command("tail", LINES_TO_READ, LOG_FILE)
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error tailing %v:%v", LOG_FILE, err)
	}

	return "\n" + string(stdout)
}

// TODO: return a data.Row
// and parse the log for errors
