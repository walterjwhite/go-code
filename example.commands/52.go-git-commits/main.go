package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	repository := os.Args[1]

	r, err := git.PlainOpen(repository)
	logging.Panic(err)

	ref, err := r.Head()

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	logging.Panic(err)

	//err = cIter.ForEach(func(c *object.Commit) error {
	//	fmt.Println(c)
	//	return nil
	//})
	//fmt.Println(cIter.Next())
	c, err := cIter.Next()
	logging.Panic(err)

	log.Info().Msgf("hash: %v", c.Hash)
	log.Info().Msg(c)
}
