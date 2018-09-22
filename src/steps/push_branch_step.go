package steps

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
	Undoable   bool
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *PushBranchStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	if step.Undoable {
		return &PushBranchAfterCurrentBranchSteps{}
	}
	return &SkipCurrentBranchSteps{}
}

// Run executes this step.
func (step *PushBranchStep) Run(deps *StepDependencies) error {
	if !deps.GitBranchService.ShouldBranchBePushed(step.BranchName) && !deps.DryRunStateService.IsActive() {
		return nil
	}
	if step.Force {
		return deps.ScriptService.RunCommand("git", "push", "-f", "origin", step.BranchName)
	}
	if deps.GitCurrentBranchService.GetCurrentBranchName() == step.BranchName {
		return deps.ScriptService.RunCommand("git", "push")
	}
	return deps.ScriptService.RunCommand("git", "push", "origin", step.BranchName)
}
