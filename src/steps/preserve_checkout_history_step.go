package steps

import (
	"github.com/Originate/git-town/src/command"
)

// PreserveCheckoutHistoryStep does stuff
type PreserveCheckoutHistoryStep struct {
	NoOpStep
	InitialBranch                     string
	InitialPreviouslyCheckedOutBranch string
}

// Run executes this step.
func (step *PreserveCheckoutHistoryStep) Run(deps *StepDependencies) error {
	expectedPreviouslyCheckedOutBranch := deps.GitBranchService.GetExpectedPreviouslyCheckedOutBranch(step.InitialPreviouslyCheckedOutBranch, step.InitialBranch)
	if expectedPreviouslyCheckedOutBranch != deps.GitBranchService.GetPreviouslyCheckedOutBranch() {
		currentBranch := deps.GitCurrentBranchService.GetCurrentBranchName()
		command.New("git", "checkout", expectedPreviouslyCheckedOutBranch).Run()
		command.New("git", "checkout", currentBranch).Run()
	}
	return nil
}
