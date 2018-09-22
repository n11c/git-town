package steps

import (
	"github.com/Originate/git-town/src/git"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *DeleteRemoteBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName}
	}
	sha := deps.GitShaService.GetBranchSha(git.GetTrackingBranchName(step.BranchName))
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: sha}
}

// Run executes this step.
func (step *DeleteRemoteBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "push", "origin", ":"+step.BranchName)
}
