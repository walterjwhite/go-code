package main

import (
	"fmt"
    "os"
)

func main() {
	fmt.Println(len(os.Args), os.Args)

	for i := 0; i < len(os.Args);i++ {
		fmt.Printf("Arg: %s\n", os.Args[i])
	}
}
