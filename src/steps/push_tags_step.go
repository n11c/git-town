package steps

// PushTagsStep pushes newly created Git tags to the remote.
type PushTagsStep struct {
	NoOpStep
}

// Run executes this step.
func (step *PushTagsStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "push", "--tags")
}
