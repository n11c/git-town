package git

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
	"go.uber.org/dig"
)

type statusService struct {
	configService      IConfigService
	environmentService IEnvironmentService
}

// NewStatusServiceOpts is the options for NewStatusService
type NewStatusServiceOpts struct {
	dig.In

	ConfigService      IConfigService
	EnvironmentService IEnvironmentService
}

// NewStatusService returns a new IStatusService
func NewStatusService(opts NewStatusServiceOpts) IStatusService {
	return &statusService{
		configService:      opts.ConfigService,
		environmentService: opts.EnvironmentService,
	}
}

// EnsureDoesNotHaveConflicts asserts that the workspace
// has no unresolved merge conflicts.
func (s *statusService) EnsureDoesNotHaveConflicts() {
	util.Ensure(!s.HasConflicts(), "You must resolve the conflicts before continuing")
}

// EnsureDoesNotHaveOpenChanges assets that the workspace
// has no open changes
func (s *statusService) EnsureDoesNotHaveOpenChanges(message string) {
	util.Ensure(!s.HasOpenChanges(), "You have uncommitted changes. "+message)
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func (s *statusService) HasConflicts() bool {
	return command.New("git", "status").OutputContainsText("Unmerged paths")
}

// HasOpenChanges returns whether the local repository contains uncommitted changes.
func (s *statusService) HasOpenChanges() bool {
	return command.New("git", "status", "--porcelain").Output() != ""
}

// HasShippableChanges returns whether the supplied branch has an changes
// not currently on the main branchName
func (s *statusService) HasShippableChanges(branchName string) bool {
	return command.New("git", "diff", s.configService.GetMainBranch()+".."+branchName).Output() != ""
}

// IsMergeInProgress returns whether the local repository is in the middle of
// an unfinished merge process.
func (s *statusService) IsMergeInProgress() bool {
	_, err := os.Stat(fmt.Sprintf("%s/.git/MERGE_HEAD", s.environmentService.GetRootDirectory()))
	return err == nil
}

// IsRebaseInProgress returns whether the local repository is in the middle of
// an unfinished rebase process.
func (s *statusService) IsRebaseInProgress() bool {
	return command.New("git", "status").OutputContainsText("rebase in progress")
}
