package logging

import (
	"log"
	"os"
)

func Set(filename string) {
	if len(filename) == 0 {
		log.Println("No log file configured, using stdout/stderr")
		return
	}

	useLogFile(filename)
}

func useLogFile(filename string) {
	log.Printf("Writing logs to: %v\n", filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Panic(err)

	// TODO: ensure the file is closed
	// defer f.Close()
	log.SetOutput(f)
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
