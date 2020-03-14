package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// An example of how to create and remove branches or any other kind of reference.
func main() {
	//CheckArgs("<url>", "<directory>")
	url, directory, commit := os.Args[1], os.Args[2], os.Args[3]

	// Clone the given repository to the given directory
	log.Info().Msgf("git clone %s %s", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url, Progress: os.Stdout,
	})
	logging.Panic(err)

	// Create a new branch to the current HEAD
	log.Info().Msg("git branch my-branch")

	ref, err := r.Head()
	logging.Panic(err)
	log.Info().Msgf("ref: %v", ref.Hash())

	w, err := r.Worktree()
	logging.Panic(err)

	logging.Panic(w.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(commit)}))
	ref, err = r.Head()
	logging.Panic(err)

	log.Info().Msgf("ref: %v", ref.Hash())
}
