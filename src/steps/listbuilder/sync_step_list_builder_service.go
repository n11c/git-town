package listbuilder

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"go.uber.org/dig"
)

type syncStepListBuilderService struct {
	gitBranchService git.IBranchService
	gitConfigService git.IConfigService
}

// NewSyncStepListBuilderServiceOpts is the options for NewSyncStepListBuilderService
type NewSyncStepListBuilderServiceOpts struct {
	dig.In

	GitBranchService git.IBranchService
	GitConfigService git.IConfigService
}

// NewSyncStepListBuilderService returns a new ISyncStepListBuilderService
func NewSyncStepListBuilderService(opts NewSyncStepListBuilderServiceOpts) ISyncStepListBuilderService {
	return &syncStepListBuilderService{
		gitBranchService: opts.GitBranchService,
		gitConfigService: opts.GitConfigService,
	}
}

// GetStepList returns the steps to sync the branch with the given name.
func (s *syncStepListBuilderService) GetStepList(branchName string, pushBranch bool) (result steps.StepList) {
	isFeature := s.gitConfigService.IsFeatureBranch(branchName)
	hasRemoteOrigin := s.gitConfigService.HasRemote("origin")

	if !hasRemoteOrigin && !isFeature {
		return
	}

	result.Append(&steps.CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		result.AppendList(s.getSyncFeatureBranchSteps(branchName))
	} else {
		result.AppendList(s.getSyncNonFeatureBranchSteps(branchName))
	}

	if pushBranch && hasRemoteOrigin && !s.gitConfigService.IsOffline() {
		if s.gitBranchService.HasTrackingBranch(branchName) {
			result.Append(&steps.PushBranchStep{BranchName: branchName})
		} else {
			result.Append(&steps.CreateTrackingBranchStep{BranchName: branchName})
		}
	}

	return
}

// Helpers

func (s *syncStepListBuilderService) getSyncFeatureBranchSteps(branchName string) (result steps.StepList) {
	if s.gitBranchService.HasTrackingBranch(branchName) {
		result.Append(&steps.MergeBranchStep{BranchName: s.gitBranchService.GetTrackingBranchName(branchName)})
	}
	result.Append(&steps.MergeBranchStep{BranchName: s.gitConfigService.GetParentBranch(branchName)})
	return
}

func (s *syncStepListBuilderService) getSyncNonFeatureBranchSteps(branchName string) (result steps.StepList) {
	if s.gitBranchService.HasTrackingBranch(branchName) {
		if s.gitConfigService.GetPullBranchStrategy() == "rebase" {
			result.Append(&steps.RebaseBranchStep{BranchName: s.gitBranchService.GetTrackingBranchName(branchName)})
		} else {
			result.Append(&steps.MergeBranchStep{BranchName: s.gitBranchService.GetTrackingBranchName(branchName)})
		}
	}

	mainBranchName := s.gitConfigService.GetMainBranch()
	if mainBranchName == branchName && s.gitConfigService.HasRemote("upstream") {
		result.Append(&steps.FetchUpstreamStep{BranchName: mainBranchName})
		result.Append(&steps.RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
	}
	return
}
