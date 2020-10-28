package runner

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"io"
	"os"
	"os/exec"
)

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

	err := command.Start()
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
	return doRun(exec.CommandContext(ctx, command, arguments...))
}

func Panic(ctx context.Context, cmd *exec.Cmd) {
	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
