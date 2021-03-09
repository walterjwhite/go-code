package main

import (
	"fmt"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
	"os"
)

func main() {

	expansions := expand.ExpandAddress(os.Args[1])

	for i := 0; i < len(expansions); i++ {
		fmt.Println(expansions[i])
	}

	p := parser.ParseAddress(os.Args[1])

	for i := 0; i < len(p); i++ {
		fmt.Printf("%s: %s\n", p[i].Label, p[i].Value)
	}

}
