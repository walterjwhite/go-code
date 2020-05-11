package comment

import (
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/timeformatter/timestamp"
	"os"
	"time"
)

const (
	commentPath        = "comments"
	commentPermissions = 0644
)

type Comment struct {
	Message  string
	DateTime time.Time

	WorkTreeConfig *git.WorkTreeConfig
}

var (
	timestampConfiguration = &timestamp.Configuration{Template: "%d" + string(os.PathSeparator) + "%d" + string(os.PathSeparator) + "%d" + string(os.PathSeparator) + "%d.%d.%d"}
)
