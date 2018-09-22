package git

import "github.com/Originate/git-town/src/command"

type shaService struct{}

// NewShaService returns a new IShaService
func NewShaService() IShaService {
	return &shaService{}
}

// GetBranchSha returns the SHA1 of the latest commit
// on the branch with the given name.
func (s *shaService) GetBranchSha(branchName string) string {
	return command.New("git", "rev-parse", branchName).Output()
}

// GetCurrentSha returns the SHA of the currently checked out commit.
func (s *shaService) GetCurrentSha() string {
	return s.GetBranchSha("HEAD")
}
