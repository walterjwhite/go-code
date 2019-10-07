package main

import (
	"fmt"
	"strings"
)

func main() {
	t1 := strings.ToLower(strings.ReplaceAll("*dnstap.Dnstap", "*", ""))

	fmt.Println(t1)
}
