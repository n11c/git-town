package steps

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoOpStep
}

// Run executes this step.
func (step *DiscardOpenChangesStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "reset", "--hard")
}
