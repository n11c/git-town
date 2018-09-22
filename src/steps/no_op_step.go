package steps

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step *NoOpStep) CreateAbortStep(deps *StepDependencies) Step {
	return &NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *NoOpStep) CreateContinueStep(deps *StepDependencies) Step {
	return &NoOpStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *NoOpStep) CreateUndoStepBeforeRun(deps *StepDependencies) Step {
	return &NoOpStep{}
}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step *NoOpStep) CreateUndoStepAfterRun(deps *StepDependencies) Step {
	return &NoOpStep{}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *NoOpStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

// Run executes this step.
func (step *NoOpStep) Run(deps *StepDependencies) error {
	return nil
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *NoOpStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
