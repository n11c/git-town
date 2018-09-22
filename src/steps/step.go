package steps

import (
	"github.com/Originate/git-town/src/dryrun"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
)

// StepDependencies
type StepDependencies struct {
	DryRunStateService              dryrun.IStateService
	GitBranchService                git.IBranchService
	GitConfigService                git.IConfigService
	GitCurrentBranchService         git.ICurrentBranchService
	GitLogService                   git.ILogService
	GitShaService                   git.IShaService
	GitSquashMergeService           git.ISquashMergeService
	GitStatusService                git.IStatusService
	GitUserService                  git.IUserService
	ScriptService                   script.IService
	SquashCommitAuthorPromptService prompt.ISquashCommitAuthorService
}

// Step represents a dedicated activity within a Git Town command.
// Git Town commands are comprised of a number of steps that need to be executed.
type Step interface {
	CreateAbortStep(deps *StepDependencies) Step
	CreateContinueStep(deps *StepDependencies) Step
	CreateUndoStepBeforeRun(deps *StepDependencies) Step
	CreateUndoStepAfterRun(deps *StepDependencies) Step
	GetAutomaticAbortErrorMessage() string
	Run(deps *StepDependencies) error
	ShouldAutomaticallyAbortOnError() bool
}
