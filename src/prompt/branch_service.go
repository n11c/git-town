package prompt

import (
	"github.com/Originate/exit"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type branchService struct{}

func NewBranchService() IBranchService {
	return &branchService{}
}

func (b *branchService) askForBranch(opts askForBranchOptions) string {
	result := ""
	prompt := &survey.Select{
		Message: opts.prompt,
		Options: opts.branchNames,
		Default: opts.defaultBranchName,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}

func (b *branchService) askForBranches(opts askForBranchesOptions) []string {
	result := []string{}
	prompt := &survey.MultiSelect{
		Message: opts.prompt,
		Options: opts.branchNames,
		Default: opts.defaultBranchNames,
	}
	err := survey.AskOne(prompt, &result, nil)
	exit.If(err)
	return result
}
