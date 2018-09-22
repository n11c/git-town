package git

// IBranchService provides methods for querying about branches
type IBranchService interface {
	DoesBranchHaveUnmergedCommits(branchName string) bool
	EnsureBranchInSync(branchName, errorMessageSuffix string)
	EnsureDoesNotHaveBranch(branchName string)
	EnsureHasBranch(branchName string)
	EnsureIsNotMainBranch(branchName, errorMessage string)
	EnsureIsNotPerennialBranch(branchName, errorMessage string)
	EnsureIsPerennialBranch(branchName, errorMessage string)
	GetBranchAuthors(branchName string) []string
	GetExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch string) string
	GetLocalBranches() []string
	GetLocalBranchesWithoutMain() []string
	GetLocalBranchesWithDeletedTrackingBranches() []string
	GetLocalBranchesWithMainBranchFirst() []string
	GetPreviouslyCheckedOutBranch() string
	GetTrackingBranchName(branchName string) string
	HasBranch(branchName string) bool
	HasLocalBranch(branchName string) bool
	HasTrackingBranch(branchName string) bool
	IsBranchInSync(branchName string) bool
	RemoveOutdatedConfiguration()
	ShouldBranchBePushed(branchName string) bool
}
