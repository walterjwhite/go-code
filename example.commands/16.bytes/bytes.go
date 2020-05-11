package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/walterjwhite/go-application/libraries/application"
)

func init() {
	application.Configure()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	keyBytes := scanner.Bytes()

	fmt.Printf("length: %v\n", len(keyBytes))
	keyBytes = append(keyBytes, '\n')

	fmt.Printf("length: %v\n", len(keyBytes))
}
