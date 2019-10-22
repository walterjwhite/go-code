package main

import (
	whois "github.com/undiabler/golang-whois"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"

	"bufio"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

const outputDirectory = "/tmp/whois/"
const timeout = 15 * time.Second

var inputFilename = flag.String("input", "", "input filename")

func main() {
	application.Configure()

	if len(*inputFilename) == 0 {
		log.Println("Input filename is required")
		os.Exit(1)
	}

	file, err := os.Open(*inputFilename)
	logging.Panic(err)

	defer file.Close()

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lookup(scanner.Text(), &wg)
	}

	err = scanner.Err()
	logging.Panic(err)

	wg.Wait()
}

func lookup(domain string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		doLookup(domain)
	}()
}

func doLookup(domain string) {
	result, err := whois.GetWhoisTimeout(domain, timeout)

	if err != nil {
		// change to Warn
		//logging.Panic(err)
		log.Printf("Error getting whois:\n%v\n", err)
		return
	}

	//fmt.Println(result)
	//fmt.Printf("Nameservers: %v \n", whois.ParseNameServers(result))

	writeWhoisResults(domain, []byte(result))
}

func writeWhoisResults(domain string, data []byte) {
	logging.Panic(os.MkdirAll(outputDirectory, 0755))

	f, err := os.Create(outputDirectory + domain)
	logging.Panic(err)

	defer f.Close()

	_, err = f.Write(data)
	logging.Panic(err)
}
