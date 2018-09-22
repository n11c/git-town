package script

// IService represents a class that run commands
type IService interface {
	ActivateDryRun()
	Fetch()
	OpenBrowser(url string)
	PrintCommand(cmd ...string)
	RunCommand(cmd ...string) error
	RunCommandSafe(cmd ...string)
	SquashMerge(branchName string)
}
