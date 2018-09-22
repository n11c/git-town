package steps

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	NoOpStep
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step *RebaseBranchStep) CreateAbortStep(deps *StepDependencies) Step {
	return &AbortRebaseBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *RebaseBranchStep) CreateContinueStep(deps *StepDependencies) Step {
	return &ContinueRebaseBranchStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *RebaseBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &ResetToShaStep{Hard: true, Sha: deps.GitShaService.GetCurrentSha()}
}

// Run executes this step.
func (step *RebaseBranchStep) Run(deps *StepDependencies) error {
	err := deps.ScriptService.RunCommand("git", "rebase", step.BranchName)
	if err != nil {
		deps.GitCurrentBranchService.ClearCurrentBranchCache()
	}
	return err
}
