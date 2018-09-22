package steps

// PullBranchStep pulls the branch with the given name from the origin remote
type PullBranchStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *PullBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "pull")
}
