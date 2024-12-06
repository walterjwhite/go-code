package jira

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func (i *Instance) Create(projectKey, summary, description, issueTypeName string) *jiral.Issue {
	i.setupAuth()

	issue := &jiral.Issue{Fields: &jiral.IssueFields{
		Project:     jiral.Project{Key: projectKey},
		Summary:     summary,
		Description: description,

		Type: jiral.IssueType{Name: issueTypeName},

	}}

	_, _, err := i.client.Issue.Create(issue)
	logging.Panic(err)

	log.Info().Msgf("Created: %v", issue.ID)

	return issue
}
