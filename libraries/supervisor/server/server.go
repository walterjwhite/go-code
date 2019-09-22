package server

import (
	//"errors"
	"log"
	"os"
	"sync"
)

const (
	Protocol   = "tcp"
	Port       = "1234"
	ListenHost = "localhost"

	SupervisionEnableFile = "/run/walterjwhite/process.supervision.enabled"
)

type Args struct {
}

type Service struct {
}

type Server string

func Status(service *Service) error {
	//update status on service object
	return nil
}

func Enable() {
	if IsEnabled() {
		return
	}

	enableFile, err := os.Create(SupervisionEnableFile)
	if err != nil {
		log.Fatal(err)
	}

	enableFile.Close()
}

func Disable() {
	if !IsEnabled() {
		return
	}

	err := os.Remove(SupervisionEnableFile)
	if err != nil {
		panic(err)
	}
}

func IsEnabled() bool {
	_, err := os.Stat(SupervisionEnableFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	return true
}

var isInitialized = false

func Refresh() {
	var wg sync.WaitGroup

	doRefresh(wg, false, RefreshBuildDateTime)
	doRefresh(wg, true, RefreshUptime)
	doRefresh(wg, true, RefreshDisks)
	doRefresh(wg, true, RefreshServices)
	doRefresh(wg, true, RefreshLogs)
	doRefresh(wg, true, RefreshInterfaces)

	// wait for all of the async methods to complete
	wg.Wait()

	isInitialized = true
}

func doRefresh(wg sync.WaitGroup, isAlwaysRefresh bool, fn func()) {
	if !isInitialized || isAlwaysRefresh {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}
}
