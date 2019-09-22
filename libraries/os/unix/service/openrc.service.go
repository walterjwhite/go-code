package service

import (
	"log"
	"os/exec"
	"strings"
)

func getServiceStatus(service *Service) string {
	cmd := exec.Command("/sbin/service", service.Name, "status")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return "Error"
	}

	status := strings.Trim(strings.Split(string(stdout), ":")[1], " \t\r\n")
	log.Printf("service (%v): %v\n", service.Name, status)

	return status
}
