package steps

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	NoOpStep
	BranchName string
	Sha        string
}

// Run executes this step.
func (step *CreateRemoteBranchStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "push", "origin", step.Sha+":refs/heads/"+step.BranchName)
}
