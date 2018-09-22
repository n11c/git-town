package steps

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

// Run executes this step.
func (step *ResetToShaStep) Run(deps *StepDependencies) error {
	if step.Sha == deps.GitShaService.GetCurrentSha() {
		return nil
	}
	cmd := []string{"git", "reset"}
	if step.Hard {
		cmd = append(cmd, "--hard")
	}
	cmd = append(cmd, step.Sha)
	return deps.ScriptService.RunCommand(cmd...)
}
