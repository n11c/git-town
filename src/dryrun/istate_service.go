package dryrun

// IStateService represents a class that holds the state of dry run
type IStateService interface {
	Activate(initialBranchName string)
	IsActive() bool
	GetCurrentBranchName() string
	SetCurrentBranchName(branchName string)
}
