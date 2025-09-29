package prompt

import (
	"context"
)

type PromptMatcherConf struct {
	MatchThreshold float64
}

func New() *PromptMatcherConf {
	return &PromptMatcherConf{MatchThreshold: 0.04}
}

type PromptMatcher interface {
	IsLocked(ctx context.Context) bool
}
