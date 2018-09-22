package steps

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *RestoreOpenChangesStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &StashOpenChangesStep{}
}

// Run executes this step.
func (step *RestoreOpenChangesStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "stash", "pop")
}
