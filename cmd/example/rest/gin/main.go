package main

import (
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go serve(&wg)

	wg.Add(1)
	go contactWorker(&wg)

	wg.Wait()
}
