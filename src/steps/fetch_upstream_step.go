package steps

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *FetchUpstreamStep) Run(deps *StepDependencies) error {
	return deps.ScriptService.RunCommand("git", "fetch", "upstream", step.BranchName)
}
