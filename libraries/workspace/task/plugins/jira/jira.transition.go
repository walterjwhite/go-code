package jira

func (j *Jira) Transition(transitionAction string) {
	j.Instance.Transition(j.IssueId, transitionAction)
}
