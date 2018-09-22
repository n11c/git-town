package git

// ICurrentBranchService provides methods around the git current branch
type ICurrentBranchService interface {
	ClearCurrentBranchCache()
	GetCurrentBranchName() string
	UpdateCurrentBranchCache(branchName string)
}
