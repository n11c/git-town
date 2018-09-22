package steps

// AddToPerennialBranches adds the branch with the given name as a perennial branch
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *AddToPerennialBranches) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &RemoveFromPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step *AddToPerennialBranches) Run(deps *StepDependencies) error {
	deps.GitConfigService.AddToPerennialBranches(step.BranchName)
	return nil
}
