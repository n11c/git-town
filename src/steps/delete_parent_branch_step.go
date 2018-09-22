package steps

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *DeleteParentBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	parent := deps.GitConfigService.GetParentBranch(step.BranchName)
	if parent == "" {
		return &NoOpStep{}
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: parent}
}

// Run executes this step.
func (step *DeleteParentBranchStep) Run(deps *StepDependencies) error {
	deps.GitConfigService.DeleteParentBranch(step.BranchName)
	return nil
}
