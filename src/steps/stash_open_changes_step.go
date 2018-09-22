package steps

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *StashOpenChangesStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &RestoreOpenChangesStep{}
}

// Run executes this step.
func (step *StashOpenChangesStep) Run(deps *StepDependencies) error {
	err := deps.ScriptService.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return deps.ScriptService.RunCommand("git", "stash")
}
