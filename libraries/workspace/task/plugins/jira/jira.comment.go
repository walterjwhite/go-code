package jira

func (j *Jira) Comment(comment string) {
	j.Instance.Comment(j.IssueId, comment)
}
