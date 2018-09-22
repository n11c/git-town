package listbuilder

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"go.uber.org/dig"
)

type appendStepListBuilderService struct {
	gitBranchService           git.IBranchService
	gitConfigService           git.IConfigService
	gitCurrentBranchService    git.ICurrentBranchService
	gitEnvironmentService      git.IEnvironmentService
	parentBranchPromptService  prompt.IParentBranchService
	scriptService              script.IService
	syncStepListBuilderService ISyncStepListBuilderService
}

type appendConfig struct {
	ParentBranch string
	TargetBranch string
}

// NewAppendStepListBuilderServiceOpts is the options for NewBranchService
type NewAppendStepListBuilderServiceOpts struct {
	dig.In

	GitBranchService           git.IBranchService
	GitConfigService           git.IConfigService
	GitCurrentBranchService    git.ICurrentBranchService
	GitEnvironmentService      git.IEnvironmentService
	ParentBranchPromptService  prompt.IParentBranchService
	ScriptService              script.IService
	SyncStepListBuilderService ISyncStepListBuilderService
}

// NewAppendStepListBuilderService returns a new IBranchService
func NewAppendStepListBuilderService(opts NewAppendStepListBuilderServiceOpts) IAppendStepListBuilderService {
	return &appendStepListBuilderService{
		gitBranchService:           opts.GitBranchService,
		gitConfigService:           opts.GitConfigService,
		gitCurrentBranchService:    opts.GitCurrentBranchService,
		gitEnvironmentService:      opts.GitEnvironmentService,
		parentBranchPromptService:  opts.ParentBranchPromptService,
		scriptService:              opts.ScriptService,
		syncStepListBuilderService: opts.SyncStepListBuilderService,
	}
}

// GetAppendStepList returns a step list for `git append`
func (a *appendStepListBuilderService) GetAppendStepList(args []string) steps.StepList {
	config := a.getAppendConfig(args)
	return a.getStepList(config)
}

// GetHackStepList returns a step list for `git hack`
func (a *appendStepListBuilderService) GetHackStepList(args []string, promptForParent bool) steps.StepList {
	config := a.getHackConfig(args, promptForParent)
	return a.getStepList(config)
}

func (a *appendStepListBuilderService) getAppendConfig(args []string) (result appendConfig) {
	result.ParentBranch = a.gitCurrentBranchService.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if a.gitConfigService.HasRemote("origin") && !a.gitConfigService.IsOffline() {
		a.scriptService.Fetch()
	}
	a.gitBranchService.EnsureDoesNotHaveBranch(result.TargetBranch)
	a.parentBranchPromptService.EnsureKnowsParentBranches([]string{result.ParentBranch})
	return
}

func (a *appendStepListBuilderService) getHackConfig(args []string, promptForParent bool) (result appendConfig) {
	result.TargetBranch = args[0]
	result.ParentBranch = a.getHackParentBranch(result.TargetBranch, promptForParent)
	if a.gitConfigService.HasRemote("origin") && !a.gitConfigService.IsOffline() {
		a.scriptService.Fetch()
	}
	a.gitBranchService.EnsureDoesNotHaveBranch(result.TargetBranch)
	return
}

func (a *appendStepListBuilderService) getHackParentBranch(targetBranch string, promptForParent bool) string {
	if promptForParent {
		parentBranch := a.parentBranchPromptService.AskForBranchParent(targetBranch, a.gitConfigService.GetMainBranch())
		a.parentBranchPromptService.EnsureKnowsParentBranches([]string{parentBranch})
		return parentBranch
	}
	return a.gitConfigService.GetMainBranch()
}

func (a *appendStepListBuilderService) getStepList(config appendConfig) (result steps.StepList) {
	for _, branchName := range append(a.gitConfigService.GetAncestorBranches(config.ParentBranch), config.ParentBranch) {
		result.AppendList(a.syncStepListBuilderService.GetStepList(branchName, true))
	}
	result.Append(&steps.CreateBranchStep{BranchName: config.TargetBranch, StartingPoint: config.ParentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.TargetBranch, ParentBranchName: config.ParentBranch})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.TargetBranch})
	if a.gitConfigService.HasRemote("origin") && a.gitConfigService.ShouldNewBranchPush() && !a.gitConfigService.IsOffline() {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.TargetBranch})
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return result
}
