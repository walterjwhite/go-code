package main

import (
	//"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"os"
	//"path/filepath"
	//"flag"
	//"syscall"
)

func init() {
	application.Configure()
}

func main() {
	/*
		chrootPath := filepath.Join("tmp", "chroot")
		logging.Panic(os.MkdirAll(chrootPath, os.ModePerm))

		log.Info().Msgf("chrooting to: %v", chrootPath)

		logging.Panic(syscall.Chroot(chrootPath))

		//_, err := runner.Run(application.Context, "echo", "hi")
		//_, err := runner.Run(application.Context, "pwd")
		command := flag.Args()[0]
		arguments := flag.Args()[1:]

		_, err := runner.Run(application.Context, command, arguments...)
		logging.Panic(err)
	*/

	cmd := runner.Prepare(application.Context /*"echo", "hi: $USER"*/, "./test")
	cmd.Env = []string{"USER=fred"}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
