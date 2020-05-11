package main

import (
	"os"

	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
)

// An example of how to create and remove branches or any other kind of reference.
func main() {
	//CheckArgs("<url>", "<directory>")
	repository := os.Args[1]

	// Clone the given repository to the given directory
	r, err := git.PlainOpen(repository)
	logging.Panic(err)

	logging.Panic(r.Push(&git.PushOptions{}))
}
