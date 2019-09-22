package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	keyBytes := scanner.Bytes()

	fmt.Printf("length: %v\n", len(keyBytes))
	keyBytes = append(keyBytes, '\n')

	fmt.Printf("length: %v\n", len(keyBytes))
}
