package steps

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortMergeBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "merge", "--abort")
}
