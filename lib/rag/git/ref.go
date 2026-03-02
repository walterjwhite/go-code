package git

import (
	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

type Ref struct {
	Kind string // "branch", "tag", "ref"
	Name string
}

func ResolveRef(repoPath, branch, tag string) Ref {
	if branch != "" {
		_, err := Output(repoPath, "rev-parse", "--verify", "refs/heads/"+branch)
		if err != nil {
			logging.Error(fmt.Errorf("branch %q not found: %w", branch, err))
		}
		return Ref{Kind: "branch", Name: branch}
	}

	if tag != "" {
		_, err := Output(repoPath, "rev-parse", "--verify", "refs/tags/"+tag)
		if err != nil {
			logging.Error(fmt.Errorf("tag %q not found: %w", tag, err))
		}
		return Ref{Kind: "tag", Name: tag}
	}

	currentBranch, err := Output(repoPath, "symbolic-ref", "--short", "-q", "HEAD")
	if err == nil && strings.TrimSpace(currentBranch) != "" {
		return Ref{Kind: "branch", Name: strings.TrimSpace(currentBranch)}
	}

	resolvedRef, resolveErr := Output(repoPath, "rev-parse", "--short", "HEAD")
	if resolveErr != nil {
		logging.Error(fmt.Errorf("resolve git ref: %w", resolveErr))
	}
	return Ref{Kind: "ref", Name: strings.TrimSpace(resolvedRef)}
}

func (r Ref) BranchOrEmpty() string {
	if r.Kind == "branch" {
		return r.Name
	}
	return ""
}

func (r Ref) TagOrEmpty() string {
	if r.Kind == "tag" {
		return r.Name
	}
	return ""
}
