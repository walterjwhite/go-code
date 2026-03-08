package runner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io"
)

var allowedCommandPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type RunConfig struct {
	Command              string
	Arguments            []string
	UseExistingEnv       bool
	EnvironmentVariables []string
	Stdout               io.Writer
	Stderr               io.Writer
}

func validateCommand(command string) error {
	if command == "" {
		return errors.New("command cannot be empty")
	}

	if strings.Contains(command, "..") {
		return errors.New("command contains invalid path traversal")
	}

	dangerousChars := []string{";", "|", "&", "$", "`", "(", ")", "{", "}", "<", ">", "\\", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(command, char) {
			return errors.New("command contains dangerous characters")
		}
	}

	if !allowedCommandPattern.MatchString(command) && !strings.Contains(command, "/") {
		return errors.New("command name contains invalid characters")
	}

	return nil
}

func validateEnvironmentVariable(envVar string) error {
	if envVar == "" {
		return errors.New("environment variable cannot be empty")
	}

	parts := strings.SplitN(envVar, "=", 2)
	if len(parts) != 2 {
		return errors.New("environment variable must be in KEY=value format")
	}

	key := parts[0]
	if key == "" {
		return errors.New("environment variable key cannot be empty")
	}

	validKey := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !validKey.MatchString(key) {
		return errors.New("environment variable key contains invalid characters")
	}

	return nil
}

func WithEnvironment(command *exec.Cmd, useExistingEnvironment bool, environmentVariables ...string) error {
	if useExistingEnvironment {
		command.Env = append(command.Env, os.Environ()...)
	}

	for _, envVar := range environmentVariables {
		if err := validateEnvironmentVariable(envVar); err != nil {
			return err
		}
		command.Env = append(command.Env, envVar)
	}

	return nil
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
	if command.Stdout == nil {
		command.Stdout = os.Stdout
	}
	if command.Stderr == nil {
		command.Stderr = os.Stderr
	}

	err := command.Start()
	if err != nil {
		log.Debug().Msg("failed to start subprocess")
		return -1, err
	}

	log.Debug().Msgf("subprocess %d", command.Process.Pid)

	err = command.Wait()
	if err != nil {
		log.Debug().Msg("subprocess exited with error")

		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		}

		return -2, err
	}

	return 0, nil
}

func Run(ctx context.Context, command string, arguments ...string) (int, error) {
	if err := validateCommand(command); err != nil {
		log.Debug().Err(err).Msg("command validation failed")
		return -1, err
	}

	for _, arg := range arguments {
		if strings.ContainsAny(arg, ";|&$`(){}<>\n\r") {
			log.Debug().Msg("argument contains dangerous characters")
			return -1, errors.New("argument contains invalid characters")
		}
	}

	return doRun(exec.CommandContext(ctx, command, arguments...))
}

func RunWithConfig(ctx context.Context, config RunConfig) (int, error) {
	if err := validateCommand(config.Command); err != nil {
		log.Debug().Err(err).Msg("command validation failed")
		return -1, err
	}

	for _, arg := range config.Arguments {
		if strings.ContainsAny(arg, ";|&$`(){}<>\n\r") {
			log.Debug().Msg("argument contains dangerous characters")
			return -1, errors.New("argument contains invalid characters")
		}
	}

	cmd := exec.CommandContext(ctx, config.Command, config.Arguments...)

	if err := WithEnvironment(cmd, config.UseExistingEnv, config.EnvironmentVariables...); err != nil {
		log.Debug().Err(err).Msg("environment validation failed")
		return -1, err
	}

	if config.Stdout != nil {
		cmd.Stdout = config.Stdout
	}
	if config.Stderr != nil {
		cmd.Stderr = config.Stderr
	}

	return doRun(cmd)
}

func Panic(ctx context.Context, cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		logging.Error(err)
		return
	}
	if err := cmd.Wait(); err != nil {
		logging.Error(err)
		return
	}
}
