package steps

// MergeBranchStep merges the branch with the given name into the current branch
type MergeBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step *MergeBranchStep) CreateAbortStep(deps *StepDependencies) Step {
	return &AbortMergeBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *MergeBranchStep) CreateContinueStep(deps *StepDependencies) Step {
	return &ContinueMergeBranchStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *MergeBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &ResetToShaStep{Hard: true, Sha: deps.GitShaService.GetCurrentSha()}
}

// Run executes this step.
func (step *MergeBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "merge", "--no-edit", step.BranchName)
}
