package service

import (
	"fmt"
	"os/exec"
	"strings"
)

func getServiceListeningAddresses(service *Service) []string {
	if len(service.Port) > 0 {
		stdout, stderr := callSs(service.Port)

		if stderr != nil {
			return []string{"Error"}
		}

		return parse(stdout)
	}

	return []string{"N/A"}
}

func callSs(port string) ([]byte, error) {
	cmd := exec.Command("/sbin/ss", "-nl", "-o", "sport", "=", fmt.Sprintf(":%v", port))
	return cmd.CombinedOutput()
}

func parse(raw []byte) []string {
	listening := make([]string, 0)

	output := strings.Split(string(raw), "\n")
	for i := 1; i < len(output); i++ {
		fields := strings.Fields(output[i])

		if len(fields) > 5 {
			listening = append(listening, fmt.Sprintf("%v:%v", fields[0], fields[4]))
		}
	}

	return listening
}
