package listbuilder

import "github.com/Originate/git-town/src/steps"

// ISyncStepListBuilderService provides a methods for getting the steps to sync a branch
type ISyncStepListBuilderService interface {
	GetStepList(branchName string, pushBranch bool) steps.StepList
}
