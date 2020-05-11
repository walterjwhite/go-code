package git

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"os/user"
	"strings"

	"bufio"
	"bytes"
	"context"
	//"errors"
	"fmt"
	"time"
)

func GetCurrentBranch(projectDirectory string) string {
	ctx, cancel := context.WithTimeout(application.Context, 5*time.Second)
	defer cancel()

	cmd := runner.Prepare( /*application.Context*/ ctx, "git", "branch")
	cmd.Dir = projectDirectory

	var b bytes.Buffer

	runner.WithWriter(cmd, &b)
	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())

	//output := string(b)

	//log.Info().Msgf("output: %v", output)

	scanner := bufio.NewScanner( /*output*/ &b)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "*") {
			return strings.TrimPrefix(line, "* ")
		}
	}
	// naming convention
	// <owner>/<source>/ticket-#
	// walterjwhite/dev/jira-123

	//logging.Panic(errors.New("Unable to parse branch (has repository been initialized?): "))
	return "master"
}

func GetOwner() string {
	currentUser, err := user.Current()
	logging.Panic(err)

	return currentUser.Username
}

func GetSourceBranch(currentBranchName string) string {
	return strings.Split(currentBranchName, "/")[1]
}

func GetTicketId(currentBranchName string) string {
	log.Info().Msgf("current branch: %v", currentBranchName)

	if strings.Contains(currentBranchName, "/") {
		return strings.Split(currentBranchName, "/")[2]
	}

	return ""
}

// TODO: generalize this
func FormatCommitMessage(projectDirectory string, messageTemplate *string, message string) string {
	currentBranchName := GetCurrentBranch(projectDirectory)
	ticketId := GetTicketId(currentBranchName)

	if len(ticketId) == 0 {
		return message
	}

	formattedMessage := fmt.Sprintf(*messageTemplate, ticketId, message)
	log.Debug().Msgf("using %v as the commit message", formattedMessage)
	return formattedMessage
}
