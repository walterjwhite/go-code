package task

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/git/plugins/submodule"

	"os"
	"path/filepath"
	"strings"
)

// updates the task status
func (t *Task) changeStatus(ctx context.Context, status string) {
	oldStatus, relativePath := t.GetStatus()

	source := filepath.Join(oldStatus, relativePath)
	target := filepath.Join(status, relativePath)

	submodule.AtomicMove(ctx, t.Workspace.WorkTreeConfig, source, target, fmt.Sprintf("status: %v -> %v (%v)", oldStatus, status, relativePath))
}

func (t *Task) GetStatus() (string, string) {
	components := strings.Split(t.Path, string(os.PathSeparator))
	status := components[0]
	relativePath := strings.Join(components[1:], string(os.PathSeparator))

	return status, relativePath
}
