package steps

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *DeleteLocalBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	sha := deps.GitShaService.GetBranchSha(step.BranchName)
	return &CreateBranchStep{BranchName: step.BranchName, StartingPoint: sha}
}

// Run executes this step.
func (step *DeleteLocalBranchStep) Run(deps *StepDependencies) error {
	op := "-d"
	if step.Force || deps.GitBranchService.DoesBranchHaveUnmergedCommits(step.BranchName) {
		op = "-D"
	}
	return deps.ScriptService.RunCommand("git", "branch", op, step.BranchName)
}
