package server

import (
	"log"
	"os/exec"
	"strings"
)

var UPTIME string

func (s *Server) Uptime(args *Args, response *string) error {
	*response = UPTIME

	log.Printf("response: %v\n", response)
	log.Printf("UPTIME: %v\n", UPTIME)
	return nil
}

func RefreshUptime() {
	cmd := exec.Command("uptime", "-p")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error updating uptime: %v", err)
		return
	}

	UPTIME = strings.Replace(string(stdout), "up", "up:", 1)
	log.Printf("Uptime: %v", UPTIME)
}
