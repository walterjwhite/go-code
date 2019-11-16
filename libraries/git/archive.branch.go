package git

import (
	"context"
	"fmt"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ArchiveRequest struct {
	BranchName string
	// do not infer the source branch
	SourceBranchName string
	BackupPath       string

	projectName   string
	patchFileName string
}

func (r *ArchiveRequest) ArchiveBranch(ctx context.Context) {
	r.getProjectName()
	r.getPatchFilename()

	_, err := runner.Run(ctx, "git", "checkout", r.SourceBranchName)
	logging.Panic(err)

	f, err := os.OpenFile(r.patchFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)
	defer f.Close()

	cmd := runner.Prepare(ctx, "git", "diff", r.SourceBranchName, r.BranchName)
	runner.WithWriter(cmd, f)
	logging.Panic(runner.Start(cmd))
	logging.Panic(runner.Wait(cmd))

	r.dropBranch(ctx)
}

func (r *ArchiveRequest) getProjectName() {
	workingDirectory, err := os.Getwd()
	logging.Panic(err)

	r.projectName = path.Base(workingDirectory)
}

func (r *ArchiveRequest) getPatchFilename() {
	workingDirectory, err := os.Getwd()
	logging.Panic(err)

	parent := filepath.Dir(workingDirectory)

	branchSafeName := strings.ReplaceAll(r.BranchName, "/", "")
	sourceBranchSafeName := strings.ReplaceAll(r.SourceBranchName, "/", "")

	r.patchFileName = filepath.Join(parent, fmt.Sprintf("%v-%v-%v.patch", r.projectName, sourceBranchSafeName, branchSafeName))
}

func (r *ArchiveRequest) dropBranch(ctx context.Context) {
	_, err := runner.Run(ctx, "git", "branch", "-D", r.BranchName)
	logging.Panic(err)
}
