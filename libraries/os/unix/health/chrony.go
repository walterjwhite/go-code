package health

import (
	"context"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ChronySources(ctx context.Context) int {
	cmd := exec.CommandContext(ctx, "chronyc", "sources", "-n")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error checking chrony sources: %v", err)
		return HEALTH_BAD
	}

	//timeout 1 chronyc sources -n | grep -P "Number of sources = [\d]{1,}" | sed -e "s/.*Number of sources = //"
	for _, line := range strings.Split(string(stdout), "\n") {
		if strings.Index(line, "Number of sources = ") > 0 {
			sourceCount := strings.Trim(strings.Split(line, "=")[1], " \t")
			i, err := strconv.Atoi(sourceCount)

			if err != nil {
				log.Printf("Error converting to int: %v / %v", sourceCount, err)
				return HEALTH_BAD
			}

			if i > 0 {
				return HEALTH_GOOD
			}

			return HEALTH_BAD
		}
	}

	return HEALTH_BAD
}
