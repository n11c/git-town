package steps

// SquashMergeBranchStep squash merges the branch with the given name into the current branch
type SquashMergeBranchStep struct {
	NoOpStep
	BranchName    string
	CommitMessage string
}

// CreateAbortStep returns the abort step for this step.
func (step *SquashMergeBranchStep) CreateAbortStep(deps *StepDependencies) Step {
	return &DiscardOpenChangesStep{}
}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step *SquashMergeBranchStep) CreateUndoStepAfterRun(deps *StepDependencies) Step {
	return &RevertCommitStep{Sha: deps.GitShaService.GetCurrentSha()}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *SquashMergeBranchStep) GetAutomaticAbortErrorMessage() string {
	return "Aborted because commit exited with error"
}

// Run executes this step.
func (step *SquashMergeBranchStep) Run(deps *StepDependencies) error {
	deps.ScriptService.SquashMerge(step.BranchName)
	commitCmd := []string{"git", "commit"}
	if step.CommitMessage != "" {
		commitCmd = append(commitCmd, "-m", step.CommitMessage)
	}
	author := deps.SquashCommitAuthorPromptService.GetSquashCommitAuthor(step.BranchName)
	if author != deps.GitUserService.GetLocalAuthor() {
		commitCmd = append(commitCmd, "--author", author)
	}
	deps.GitSquashMergeService.CommentOutSquashCommitMessage("")
	return deps.ScriptService.RunCommand(commitCmd...)
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
