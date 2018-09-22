package steps

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// CommitOpenChangesStep commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *CommitOpenChangesStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	branchName := git.GetCurrentBranchName()
	return &ResetToShaStep{Sha: deps.GitShaService.GetBranchSha(branchName)}
}

// Run executes this step.
func (step *CommitOpenChangesStep) Run(deps *StepDependencies) error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "commit", "-m", fmt.Sprintf("WIP on %s", git.GetCurrentBranchName()))
}
