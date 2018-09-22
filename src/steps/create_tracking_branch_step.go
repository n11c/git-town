package steps

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *CreateTrackingBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &DeleteRemoteBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step *CreateTrackingBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "push", "-u", "origin", step.BranchName)
}
