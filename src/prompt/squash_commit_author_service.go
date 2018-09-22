package prompt

import (
	"fmt"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/git"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type squashCommitAuthorService struct {
	gitBranchService git.IBranchService
}

// NewSquashCommitAuthorService returns a new ISquashCommitAuthorService
func NewSquashCommitAuthorService(gitBranchService git.IBranchService) ISquashCommitAuthorService {
	return &squashCommitAuthorService{
		gitBranchService: gitBranchService,
	}
}

// GetSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func (s *squashCommitAuthorService) GetSquashCommitAuthor(branchName string) string {
	authors := s.gitBranchService.GetBranchAuthors(branchName)
	if len(authors) == 1 {
		return authors[0]
	}
	cfmt.Printf(squashCommitAuthorHeaderTemplate, branchName)
	fmt.Println()
	return s.askForAuthor(authors)
}

// Helpers

func (s *squashCommitAuthorService) askForAuthor(authors []string) string {
	result := ""
	prompt := &survey.Select{
		Message: "Please choose an author for the squash commit:",
		Options: authors,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}
