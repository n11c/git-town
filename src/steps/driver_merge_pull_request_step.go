package steps

import (
	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/drivers"
)

// DriverMergePullRequestStep squash merges the branch with the given name into the current branch
type DriverMergePullRequestStep struct {
	NoOpStep
	BranchName                string
	CommitMessage             string
	DefaultCommitMessage      string
	enteredEmptyCommitMessage bool
	mergeError                error
	mergeSha                  string
}

// CreateAbortStep returns the abort step for this step.
func (step *DriverMergePullRequestStep) CreateAbortStep(deps *StepDependencies) Step {
	if step.enteredEmptyCommitMessage {
		return &DiscardOpenChangesStep{}
	}
	return nil
}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step *DriverMergePullRequestStep) CreateUndoStepAfterRun(deps *StepDependencies) Step {
	return &RevertCommitStep{Sha: step.mergeSha}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *DriverMergePullRequestStep) GetAutomaticAbortErrorMessage() string {
	if step.enteredEmptyCommitMessage {
		return "Aborted because commit exited with error"
	}
	return step.mergeError.Error()
}

// Run executes this step.
func (step *DriverMergePullRequestStep) Run(deps *StepDependencies) error {
	commitMessage := step.CommitMessage
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a driver
		// then revert the commit since merging via the driver will perform the actual squash merge
		step.enteredEmptyCommitMessage = true
		deps.ScriptService.SquashMerge(step.BranchName)
		deps.GitSquashMergeService.CommentOutSquashCommitMessage(step.DefaultCommitMessage + "\n\n")
		err := deps.ScriptService.RunCommand("git", "commit")
		if err != nil {
			return err
		}
		commitMessage = deps.GitLogService.GetLastCommitMessage()
		err = deps.ScriptService.RunCommand("git", "reset", "--hard", "HEAD~1")
		exit.IfWrap(err, "Error resetting the main branch")
		step.enteredEmptyCommitMessage = false
	}
	driver := drivers.GetActiveDriver()
	step.mergeSha, step.mergeError = driver.MergePullRequest(drivers.MergePullRequestOptions{
		Branch:        step.BranchName,
		CommitMessage: commitMessage,
		LogRequests:   true,
		ParentBranch:  deps.GitCurrentBranchService.GetCurrentBranchName(),
	})
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *DriverMergePullRequestStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
