package git

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
	"go.uber.org/dig"
)

type branchService struct {
	configService             IConfigService
	remoteBranches            []string
	remoteBranchesInitialized bool
	shaService                IShaService
}

// NewBranchServiceOpts is the options for NewBranchService
type NewBranchServiceOpts struct {
	dig.In

	ConfigService IConfigService
	ShaService    IShaService
}

// NewBranchService returns a new IBranchService
func NewBranchService(opts NewBranchServiceOpts) IBranchService {
	return &branchService{
		configService: opts.ConfigService,
		shaService:    opts.ShaService,
	}
}

// DoesBranchHaveUnmergedCommits returns whether the branch with the given name
// contains commits that are not merged into the main branch
func (b *branchService) DoesBranchHaveUnmergedCommits(branchName string) bool {
	return command.New("git", "log", b.configService.GetMainBranch()+".."+branchName).Output() != ""
}

// EnsureBranchInSync enforces that a branch with the given name is in sync with its tracking branch
func (b *branchService) EnsureBranchInSync(branchName, errorMessageSuffix string) {
	util.Ensure(b.IsBranchInSync(branchName), fmt.Sprintf("'%s' is not in sync with its tracking branch. %s", branchName, errorMessageSuffix))
}

// EnsureDoesNotHaveBranch enforces that a branch with the given name does not exist
func (b *branchService) EnsureDoesNotHaveBranch(branchName string) {
	util.Ensure(!b.HasBranch(branchName), fmt.Sprintf("A branch named '%s' already exists", branchName))
}

// EnsureHasBranch enforces that a branch with the given name exists
func (b *branchService) EnsureHasBranch(branchName string) {
	util.Ensure(b.HasBranch(branchName), fmt.Sprintf("There is no branch named '%s'", branchName))
}

// EnsureIsNotMainBranch enforces that a branch with the given name is not the main branch
func (b *branchService) EnsureIsNotMainBranch(branchName, errorMessage string) {
	util.Ensure(!b.configService.IsMainBranch(branchName), errorMessage)
}

// EnsureIsNotPerennialBranch enforces that a branch with the given name is not a perennial branch
func (b *branchService) EnsureIsNotPerennialBranch(branchName, errorMessage string) {
	util.Ensure(!b.configService.IsPerennialBranch(branchName), errorMessage)
}

// EnsureIsPerennialBranch enforces that a branch with the given name is a perennial branch
func (b *branchService) EnsureIsPerennialBranch(branchName, errorMessage string) {
	util.Ensure(b.configService.IsPerennialBranch(branchName), errorMessage)
}

// GetBranchAuthors returns the authors of the branch
func (b *branchService) GetBranchAuthors(branchName string) (result []string) {
	// Returns lines of "<number of commits>\t<name and email>"
	output := command.New("git", "shortlog", "-s", "-n", "-e", b.configService.GetMainBranch()+".."+branchName).Output()
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return
}

// GetExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs
func (b *branchService) GetExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch string) string {
	if b.HasLocalBranch(initialPreviouslyCheckedOutBranch) {
		if GetCurrentBranchName() == initialBranch || !b.HasLocalBranch(initialBranch) {
			return initialPreviouslyCheckedOutBranch
		}
		return initialBranch
	}
	return b.configService.GetMainBranch()
}

// GetLocalBranches returns the names of all branches in the local repository,
// ordered alphabetically
func (b *branchService) GetLocalBranches() (result []string) {
	for _, line := range strings.Split(command.New("git", "branch").Output(), "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

// GetLocalBranchesWithoutMain returns the names of all branches in the local repository,
// ordered alphabetically without the main branch
func (b *branchService) GetLocalBranchesWithoutMain() (result []string) {
	mainBranch := b.configService.GetMainBranch()
	for _, branch := range b.GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetLocalBranchesWithDeletedTrackingBranches returns the names of all branches
// whose remote tracking branches have been deleted
func (b *branchService) GetLocalBranchesWithDeletedTrackingBranches() (result []string) {
	for _, line := range strings.Split(command.New("git", "branch", "-vv").Output(), "\n") {
		line = strings.Trim(line, "* ")
		parts := strings.SplitN(line, " ", 2)
		branchName := parts[0]
		deleteTrackingBranchStatus := fmt.Sprintf("[%s: gone]", b.GetTrackingBranchName(branchName))
		if strings.Contains(parts[1], deleteTrackingBranchStatus) {
			result = append(result, branchName)
		}
	}
	return
}

// GetLocalBranchesWithMainBranchFirst returns the names of all branches
// that exist in the local repository,
// ordered to have the name of the main branch first,
// then the names of the branches, ordered alphabetically
func (b *branchService) GetLocalBranchesWithMainBranchFirst() (result []string) {
	mainBranch := GetMainBranch()
	result = append(result, mainBranch)
	for _, branch := range b.GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetPreviouslyCheckedOutBranch returns the name of the previously checked out branch
func (b *branchService) GetPreviouslyCheckedOutBranch() string {
	cmd := command.New("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if cmd.Err() != nil {
		return ""
	}
	return cmd.Output()
}

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name
func (b *branchService) GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

// HasBranch returns whether the repository contains a branch with the given name.
// The branch does not have to be present on the local repository.
func (b *branchService) HasBranch(branchName string) bool {
	for _, line := range strings.Split(command.New("git", "branch", "-a").Output(), "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "remotes/origin/", "", 1)
		if line == branchName {
			return true
		}
	}
	return false
}

// HasLocalBranch returns whether the local repository contains
// a branch with the given name.
func (b *branchService) HasLocalBranch(branchName string) bool {
	return util.DoesStringArrayContain(b.GetLocalBranches(), branchName)
}

// HasTrackingBranch returns whether the local branch with the given name
// has a tracking branch.
func (b *branchService) HasTrackingBranch(branchName string) bool {
	trackingBranchName := b.GetTrackingBranchName(branchName)
	for _, line := range b.getRemoteBranches() {
		if strings.TrimSpace(line) == trackingBranchName {
			return true
		}
	}
	return false
}

// IsBranchInSync returns whether the branch with the given name is in sync with its tracking branch
func (b *branchService) IsBranchInSync(branchName string) bool {
	if b.HasTrackingBranch(branchName) {
		localSha := b.shaService.GetBranchSha(branchName)
		remoteSha := b.shaService.GetBranchSha(b.GetTrackingBranchName(branchName))
		return localSha == remoteSha
	}
	return true
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration
func (b *branchService) RemoveOutdatedConfiguration() {
	for child, parent := range b.configService.GetParentBranchMap() {
		if !b.HasBranch(child) || !b.HasBranch(parent) {
			b.configService.DeleteParentBranch(child)
		}
	}
}

// ShouldBranchBePushed returns whether the local branch with the given name
// contains commits that have not been pushed to the remote.
func (b *branchService) ShouldBranchBePushed(branchName string) bool {
	trackingBranchName := b.GetTrackingBranchName(branchName)
	cmd := command.New("git", "rev-list", "--left-right", branchName+"..."+trackingBranchName)
	return cmd.Output() != ""
}

// Helpers

func (b *branchService) getRemoteBranches() []string {
	if !b.remoteBranchesInitialized {
		b.remoteBranches = strings.Split(command.New("git", "branch", "-r").Output(), "\n")
		b.remoteBranchesInitialized = true
	}
	return remoteBranches
}
