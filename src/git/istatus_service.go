package git

// IStatusService provides methods around conflicts and changes
type IStatusService interface {
	EnsureDoesNotHaveConflicts()
	EnsureDoesNotHaveOpenChanges(message string)
	HasConflicts() bool
	HasOpenChanges() bool
	HasShippableChanges(branchName string) bool
	IsMergeInProgress() bool
	IsRebaseInProgress() bool
}
