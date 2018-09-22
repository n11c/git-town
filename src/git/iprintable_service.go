package git

// IPrintableService provides methods for printing Git Town configuration
type IPrintableService interface {
	GetPrintableMainBranch() string
	GetPrintablePerennialBranches() string
	GetPrintablePerennialBranchTrees() string
	GetPrintableNewBranchPushFlag()
	GetPrintableBranchTree(branchName string) (result string)
	GetPrintableOfflineFlag()
}
