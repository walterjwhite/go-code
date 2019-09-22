package service

import (
	"../health"
	"log"
	"os/exec"
)

func repair(service *Service) {
	if service.Status != "started" || service.Health == health.HEALTH_BAD {
		log.Printf("Status: %v\n", service.Status)
		runRestart(service)

		// check if repaired
	}
}

func runRestart(service *Service) bool {
	log.Printf("Restarting service: %v\n", service.Name)
	if len(service.Restart) > 0 {
		for _, restartCommand := range service.Restart {
			if !run(restartCommand) {
				log.Printf("Running restart command (%v) for service service: %v failed\n", service.Restart, service.Name)
				return false
			}
		}

		log.Printf("Running restart command (%v) for service service: %v completed\n", service.Restart, service.Name)

		return true
	}

	if !serviceAction(service, "restart") {
		log.Printf("Service restart: %v failed\n", service.Name)

		if !zap(service) {
			log.Printf("Service zap: %v failed\n", service.Name)

			return kill(service)
		}
	}

	log.Printf("Successfully restarted service: %v\n", service.Name)

	return true
}

func zap(service *Service) bool {
	serviceAction(service, "zap")

	return serviceAction(service, "restart")
}

func kill(service *Service) bool {
	run("killall", service.Name)

	return serviceAction(service, "restart")
}

func serviceAction(service *Service, action string) bool {
	return run("service", service.Name, action)
}

func run(command string, arguments ...string) bool {
	cmd := exec.Command(command, arguments...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error running restart command: %v / %v", command, err)
		return false
	}

	return true
}
