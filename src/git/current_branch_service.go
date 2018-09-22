package git

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/dryrun"
	"go.uber.org/dig"
)

// CurrentBranchService provides methods around the git current branch
type CurrentBranchService struct {
	currentBranchCache string // cached to minimize the number of git commands run
	dryrunStateService dryrun.IStateService
	statusService      IStatusService
}

type NewCurrentBranchServiceOpts struct {
	dig.In

	DryrunStateService dryrun.IStateService
	StatusService      IStatusService
}

func NewCurrentBranchService(opts NewCurrentBranchServiceOpts) ICurrentBranchService {
	return &CurrentBranchService{
		dryrunStateService: opts.DryrunStateService,
		statusService:      opts.StatusService,
	}
}

// GetCurrentBranchName returns the name of the currently checked out branch
func (c *CurrentBranchService) GetCurrentBranchName() string {
	if c.dryrunStateService.IsActive() {
		return c.dryrunStateService.GetCurrentBranchName()
	}
	if c.currentBranchCache == "" {
		if c.statusService.IsRebaseInProgress() {
			currentBranchCache = c.getCurrentBranchNameDuringRebase()
		} else {
			currentBranchCache = command.New("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		}
	}
	return currentBranchCache
}

// ClearCurrentBranchCache clears the cache of the current branch.
// This should be called when a rebase fails
func (c *CurrentBranchService) ClearCurrentBranchCache() {
	c.currentBranchCache = ""
}

// UpdateCurrentBranchCache clears the cache of the current branch.
// This should be called when a new branch is checked out
func (c *CurrentBranchService) UpdateCurrentBranchCache(branchName string) {
	c.currentBranchCache = branchName
}

// Helpers

func (c *CurrentBranchService) getCurrentBranchNameDuringRebase() string {
	filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
	rawContent, err := ioutil.ReadFile(filename)
	exit.If(err)
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
