package listbuilder

import (
	"github.com/Originate/git-town/src/steps"
)

// IAppendStepListBuilderService provides methods for creating step lists for append and hack commands
type IAppendStepListBuilderService interface {
	GetAppendStepList(args []string) steps.StepList
	GetHackStepList(args []string, promptForParent bool) steps.StepList
}
