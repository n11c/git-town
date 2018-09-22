package prompt

// IParentBranchService provides methods for prompting for parent branches
type IParentBranchService interface {
	EnsureKnowsParentBranches(branchNames []string)
	AskForBranchAncestry(branchName, defaultBranchName string)
	AskForBranchParent(branchName, defaultBranchName string) string
}
