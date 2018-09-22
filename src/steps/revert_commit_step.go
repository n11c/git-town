package steps

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoOpStep
	Sha string
}

// Run executes this step.
func (step *RevertCommitStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "revert", step.Sha)
}
