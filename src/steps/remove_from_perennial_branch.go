package steps

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *RemoveFromPerennialBranches) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &AddToPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step *RemoveFromPerennialBranches) Run(deps *StepDependencies) error {
	deps.GitConfigService.RemoveFromPerennialBranches(step.BranchName)
	return nil
}
