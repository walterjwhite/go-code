package runner

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
)

func Prepare(ctx context.Context, command string, arguments ...string) *exec.Cmd {
	log.Debug().Msgf("running %v %v with %v", command, arguments, ctx)

	return exec.CommandContext(ctx, command, arguments...)
}

func WithEnvironment(command *exec.Cmd, useExistingEnvironment bool, environmentVariables ...string) {
	if useExistingEnvironment {
		command.Env = append(command.Env, os.Environ()...)
	}

	command.Env = append(command.Env, environmentVariables...)
}

func WithWriter(command *exec.Cmd, writer io.Writer) {
	command.Stdout = writer
	command.Stderr = writer
}

func WithWriters(command *exec.Cmd, writers ...io.Writer) {
	command.Stdout = io.MultiWriter(writers...)
	command.Stderr = io.MultiWriter(writers...)
}

func doRun(command *exec.Cmd) (int, error) {
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	err := Start(command)
	if err != nil {
		return -1, err
	}

	log.Debug().Msgf("subprocess %d", command.Process.Pid)

	err = command.Wait()
	if err != nil {
		log.Print(err)

		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		}

		return -2, err
	}

	return 0, nil
}

func Run(ctx context.Context, command string, arguments ...string) (int, error) {
	return doRun(Prepare(ctx, command, arguments...))
}

func Start(command *exec.Cmd) error {
	err := command.Start()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func Wait(command *exec.Cmd) error {
	err := command.Wait()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
