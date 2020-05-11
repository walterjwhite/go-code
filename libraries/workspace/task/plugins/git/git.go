package git

type Git struct {
	Url string

	// used for initial checkout and updates
	SourceBranch string

	// used for final push
	TargetBranch string

	// merge to target branch on close
	MergeOnCompletion bool

	OpenMergeRequestOnCompletion bool
}
