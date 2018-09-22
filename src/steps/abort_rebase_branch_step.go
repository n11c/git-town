package steps

// AbortRebaseBranchStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortRebaseBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "rebase", "--abort")
}
