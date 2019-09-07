package logging

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
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}

	// defer f.Close()
	log.SetOutput(f)
}
