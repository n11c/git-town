package steps

import (
	"github.com/Originate/git-town/src/drivers"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *CreatePullRequestStep) Run(deps *StepDependencies) error {
	driver := drivers.GetActiveDriver()
	parentBranch := deps.GitConfigService.GetParentBranch(step.BranchName)
	deps.ScriptService.OpenBrowser(driver.GetNewPullRequestURL(step.BranchName, parentBranch))
	return nil
}
