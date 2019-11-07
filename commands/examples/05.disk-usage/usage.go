package main

import (
	"flag"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/io/disk"

	"log"
	"strings"
)

var (
	mountPoints = flag.String("MountPoints", "/", "MountPoints to query")
)

func init() {
	application.Configure()
}

func main() {
	for _, mountPoint := range getMountPoints(*mountPoints) {
		diskUsage := disk.Usage(mountPoint)
		log.Printf("%v disk: %v / %v\n", diskUsage.MountPoint, diskUsage.UsagePercentage, diskUsage.Free)
	}
}

func getMountPoints(mountPointsString string) []string {
	return strings.Split(mountPointsString, ",")
}
