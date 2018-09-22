package steps

import (
	"os"

	"github.com/Originate/exit"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	NoOpStep
	Directory string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *ChangeDirectoryStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	dir, err := os.Getwd()
	exit.If(err)
	return &ChangeDirectoryStep{Directory: dir}
}

// Run executes this step.
func (step *ChangeDirectoryStep) Run(deps *StepDependencies) error {
	_, err := os.Stat(step.Directory)
	if err == nil {
		deps.ScriptService.PrintCommand("cd", step.Directory)
		return os.Chdir(step.Directory)
	}
	return nil
}
