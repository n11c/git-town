package git

// IShaService provides methods for getting shas of branch
type IShaService interface {
	GetBranchSha(branchName string) string
	GetCurrentSha() string
}
