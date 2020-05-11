package comment

import (
	"github.com/walterjwhite/go-application/libraries/git"
	"time"
)

func New(w *git.WorkTreeConfig, message string) *Comment {
	return &Comment{Message: message, DateTime: time.Now(), WorkTreeConfig: w}
}
