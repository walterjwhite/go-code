package server

import (
	"github.com/walterjwhite/go-application/libraries/io/disk"
	"log"
)

var disks []disk.Disk

type DiskServer []disk.Disk

func (s *DiskServer) Disks(args *Args, response *[]disk.Disk) error {
	*response = disks
	return nil
}

func RefreshDisks() {
	disks = disk.Data()
	log.Printf("disks: %v", disks)
}
