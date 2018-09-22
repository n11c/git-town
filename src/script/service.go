package script

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/browsers"
	"github.com/Originate/git-town/src/dryrun"
	"github.com/Originate/git-town/src/git"
	"go.uber.org/dig"

	"github.com/fatih/color"
)

type service struct {
	browserOpenService browsers.IOpenService
	dryRunStateService dryrun.IStateService
}

// ServiceOpts is the options for NewService
type ServiceOpts struct {
	dig.In

	BrowserOpenService browsers.IOpenService
	DryRunStateService dryrun.IStateService
}

// NewService returns a new IService
func NewService(opts ServiceOpts) IService {
	return &service{
		browserOpenService: opts.BrowserOpenService,
		dryRunStateService: opts.DryRunStateService,
	}
}

// ActivateDryRun causes all commands to not be run
func (s *service) ActivateDryRun() {
	_, err := color.New(color.FgBlue).Print(dryRunMessage)
	exit.If(err)
	s.dryRunStateService.Activate(git.GetCurrentBranchName())
}

func (s *service) Fetch() {
	err := s.RunCommand("git", "fetch", "--prune", "--tags")
	exit.If(err)
}

// OpenBrowser opens the default browser with the given URL.
func (s *service) OpenBrowser(url string) {
	command := s.browserOpenService.GetOpenBrowserCommand()
	err := RunCommand(command, url)
	exit.If(err)
}

// PrintCommand prints the given command-line operation on the console.
func (s *service) PrintCommand(cmd ...string) {
	header := ""
	for index, part := range cmd {
		if strings.Contains(part, " ") {
			part = "\"" + strings.Replace(part, "\"", "\\\"", -1) + "\""
		}
		if index != 0 {
			header = header + " "
		}
		header = header + part
	}
	if strings.HasPrefix(header, "git") && git.IsRepository() {
		header = fmt.Sprintf("[%s] %s", git.GetCurrentBranchName(), header)
	}
	fmt.Println()
	_, err := color.New(color.Bold).Println(header)
	exit.If(err)
}

// RunCommand executes the given command-line operation.
func (s *service) RunCommand(cmd ...string) error {
	PrintCommand(cmd...)
	if s.dryRunStateService.IsActive() {
		if len(cmd) == 3 && cmd[0] == "git" && cmd[1] == "checkout" {
			s.dryRunStateService.SetCurrentBranchName(cmd[2])
		}
		return nil
	}
	// Windows commands run inside CMD
	// because opening browsers is done via "start"
	if runtime.GOOS == "windows" {
		cmd = append([]string{"cmd", "/C"}, cmd...)
	}
	subProcess := exec.Command(cmd[0], cmd[1:]...) // #nosec
	subProcess.Stderr = os.Stderr
	subProcess.Stdin = os.Stdin
	subProcess.Stdout = os.Stdout
	return subProcess.Run()
}

// RunCommandSafe executes the given command-line operation, exiting if the command errors
func (s *service) RunCommandSafe(cmd ...string) {
	err := RunCommand(cmd...)
	exit.If(err)
}

// SquashMerge squash merges the given branch into the current branch
func (s *service) SquashMerge(branchName string) {
	err := s.RunCommand("git", "merge", "--squash", branchName)
	exit.IfWrap(err, "Error squash merging")
}
