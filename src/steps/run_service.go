package steps

import (
	"fmt"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/dryrun"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/util"
	"go.uber.org/dig"

	"github.com/fatih/color"
)

type RunService struct {
	runStateToDiskService IRunStateToDiskService
	stepDependencies      *StepDependencies
}

type NewRunServiceOpts struct {
	dig.In

	DryRunStateService              dryrun.IStateService
	GitBranchService                git.IBranchService
	GitConfigService                git.IConfigService
	GitCurrentBranchService         git.ICurrentBranchService
	GitLogService                   git.ILogService
	GitShaService                   git.IShaService
	GitSquashMergeService           git.ISquashMergeService
	GitStatusService                git.IStatusService
	GitUserService                  git.IUserService
	RunStateToDiskService           IRunStateToDiskService
	ScriptService                   script.IService
	SquashCommitAuthorPromptService prompt.ISquashCommitAuthorService
}

func NewRunService(opts NewRunServiceOpts) IRunService {
	return &RunService{
		runStateToDiskService: opts.RunStateToDiskService,
		stepDependencies: &StepDependencies{
			DryRunStateService:              opts.DryRunStateService,
			GitBranchService:                opts.GitBranchService,
			GitConfigService:                opts.GitConfigService,
			GitCurrentBranchService:         opts.GitCurrentBranchService,
			GitLogService:                   opts.GitLogService,
			GitShaService:                   opts.GitShaService,
			GitSquashMergeService:           opts.GitSquashMergeService,
			GitStatusService:                opts.GitStatusService,
			GitUserService:                  opts.GitUserService,
			ScriptService:                   opts.ScriptService,
			SquashCommitAuthorPromptService: opts.SquashCommitAuthorPromptService,
		},
	}
}

// Run runs the Git Town command described by the given state
// nolint: gocyclo
func (r *RunService) Run(runState *RunState) {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if !runState.IsAbort && !runState.isUndo {
				SaveRunState(runState)
			}
			fmt.Println()
			return
		}
		if getTypeName(step) == "*SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if getTypeName(step) == "*PushBranchAfterCurrentBranchSteps" {
			runState.AddPushBranchStepAfterCurrentBranchSteps()
			continue
		}
		undoStepBeforeRun := step.CreateUndoStepBeforeRun(r.stepDependencies)
		err := step.Run(r.stepDependencies)
		if err != nil {
			runState.AbortStepList.Append(step.CreateAbortStep(r.stepDependencies))
			if step.ShouldAutomaticallyAbortOnError() {
				abortRunState := runState.CreateAbortRunState()
				r.Run(&abortRunState)
				util.ExitWithErrorMessage(step.GetAutomaticAbortErrorMessage())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep(r.stepDependencies))
				runState.MarkAsUnfinished()
				if runState.Command == "sync" && !(git.IsRebaseInProgress() && git.IsMainBranch(git.GetCurrentBranchName())) {
					runState.UnfinishedDetails.CanSkip = true
				}
				r.runStateToDiskService.SaveRunState(runState)
				r.exitWithMessages(runState.UnfinishedDetails.CanSkip)
			}
		}
		undoStepAfterRun := step.CreateUndoStepAfterRun(r.stepDependencies)
		runState.UndoStepList.Prepend(undoStepBeforeRun)
		runState.UndoStepList.Prepend(undoStepAfterRun)
	}
}

// Helpers

func (r *RunService) exitWithMessages(canSkip bool) {
	messageFmt := color.New(color.FgRed)
	fmt.Println()
	_, err := messageFmt.Printf("To abort, run \"git-town abort\".\n")
	exit.If(err)
	_, err = messageFmt.Printf("To continue after having resolved conflicts, run \"git-town continue\".\n")
	exit.If(err)
	if canSkip {
		_, err = messageFmt.Printf("To continue by skipping the current branch, run \"git-town skip\".\n")
		exit.If(err)
	}
	fmt.Println()
	os.Exit(1)
}
