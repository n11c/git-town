package prompt

type IBranchService interface {
	askForBranch(opts askForBranchOptions) string
	askForBranches(opts askForBranchesOptions) []string
}
