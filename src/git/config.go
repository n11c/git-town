/*
This file contains functionality around storing configuration settings
inside Git's metadata storage for the repository.
*/

package git

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// Configuration manages the Git Town configuration
// stored in Git metadata in the given local repo and the global Git configuration.
// This class manages which config values are stored in local vs global settings.
type Configuration struct {

	// localDir contains the directory of the local Git repo.
	localDir string

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache map[string]string

	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache map[string]string
}

// Config provides the current configuration.
// This is used in the Git Town business logic, which runs in the current directory.
// The configuration is lazy-loaded this way to allow using some Git Town commands outside of Git repositories.
func Config() *Configuration {
	if currentDirConfig == nil {
		currentDirConfig = NewConfiguration("")
	}
	return currentDirConfig
}

// currentDirConfig contains the Git Town configuration in the current working directory.
var currentDirConfig *Configuration

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewConfiguration(dir string) *Configuration {
	return &Configuration{
		localDir:          dir,
		localConfigCache:  loadGitConfig(dir, false),
		globalConfigCache: loadGitConfig(dir, true),
	}
}

// loadGitConfig provides the Git configuration from the given directory or the global one if the global flag is set.
func loadGitConfig(dir string, global bool) map[string]string {
	result := map[string]string{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	res := command.RunInDir(dir, "git", cmdArgs...)
	if res.Err() != nil {
		if strings.Contains(res.OutputSanitized(), "No such file or directory") {
			return result
		}
		panic(res.Err())
	}
	output := res.Output()
	if output == "" {
		return result
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		result[key] = value
	}
	return result
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (c *Configuration) AddToPerennialBranches(branchNames ...string) {
	c.SetPerennialBranches(append(c.GetPerennialBranches(), branchNames...))
}

// AddGitAlias sets the given Git alias.
func (c *Configuration) AddGitAlias(command string) *command.Result {
	return c.setGlobalConfigValue("alias."+command, "town "+command)
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Configuration) DeleteParentBranch(branchName string) {
	c.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// EnsureIsFeatureBranch asserts that the given branch is a feature branch.
func (c *Configuration) EnsureIsFeatureBranch(branchName, errorSuffix string) {
	util.Ensure(c.IsFeatureBranch(branchName), fmt.Sprintf("The branch '%s' is not a feature branch. %s", branchName, errorSuffix))
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (c *Configuration) GetAncestorBranches(branchName string) (result []string) {
	parentBranchMap := c.GetParentBranchMap()
	current := branchName
	for {
		if c.IsMainBranch(current) || c.IsPerennialBranch(current) {
			return
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func (c *Configuration) GetChildBranches(branchName string) (result []string) {
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := c.getLocalConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetCodeHostingDriverName provides the name of the code hosting driver to use.
func (c *Configuration) GetCodeHostingDriverName() string {
	return c.getLocalOrGlobalConfigValue("git-town.code-hosting-driver")
}

// GetCodeHostingOriginHostname provides the host name of the code hosting server.
func (c *Configuration) GetCodeHostingOriginHostname() string {
	return c.getLocalConfigValue("git-town.code-hosting-origin-hostname")
}

// getGlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getGlobalConfigValue(key string) string {
	return c.globalConfigCache[key]
}

// getLocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getLocalConfigValue(key string) string {
	return c.localConfigCache[key]
}

// getGlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getLocalOrGlobalConfigValue(key string) string {
	local := c.getLocalConfigValue(key)
	if local != "" {
		return local
	}
	return c.getGlobalConfigValue(key)
}

// GetParentBranchMap returns a map from branch name to its parent branch
func (c *Configuration) GetParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.getLocalConfigValue(key)
		result[child] = parent
	}
	return result
}

// GetGitAlias provides the currently set alias for the given Git Town command.
func (c *Configuration) GetGitAlias(command string) string {
	return c.getGlobalConfigValue("alias." + command)
}

// GetGitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (c *Configuration) GetGitHubToken() string {
	return c.getLocalOrGlobalConfigValue("git-town.github-token")
}

// GetMainBranch returns the name of the main branch.
func (c *Configuration) GetMainBranch() string {
	return c.getLocalOrGlobalConfigValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func (c *Configuration) GetParentBranch(branchName string) string {
	return c.getLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func (c *Configuration) GetPerennialBranches() []string {
	result := c.getLocalOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func (c *Configuration) GetPullBranchStrategy() string {
	config := c.getLocalOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Configuration) GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := c.getLocalConfigValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.RunInDir(c.localDir, "git", "remote", "get-url", "origin").OutputSanitized()
}

// GetURLHostname returns the hostname contained within the given Git URL.
func (c *Configuration) GetURLHostname(url string) string {
	hostnameRegex := regexp.MustCompile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func (c *Configuration) GetURLRepositoryName(url string) string {
	hostname := c.GetURLHostname(url)
	repositoryNameRegex := regexp.MustCompile(".*" + hostname + "[/:](.+)")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}

// HasParentBranch returns whether or not the given branch has a parent
func (c *Configuration) HasParentBranch(branchName string) bool {
	return c.GetParentBranch(branchName) != ""
}

// IsAncestorBranch indicates whether the given branch is an ancestor of the other given branch.
func (c *Configuration) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := c.GetAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (c *Configuration) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.IsPerennialBranch(branchName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (c *Configuration) IsMainBranch(branchName string) bool {
	return branchName == c.GetMainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode
func (c *Configuration) IsOffline() bool {
	config := c.getGlobalConfigValue("git-town.offline")
	if config != "" {
		return util.StringToBool(config)
	}
	return false
}

// IsPerennialBranch indicates whether the branch with the given name is
// a perennial branch.
func (c *Configuration) IsPerennialBranch(branchName string) bool {
	perennialBranches := c.GetPerennialBranches()
	return util.DoesStringArrayContain(perennialBranches, branchName)
}

// localConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (c *Configuration) localConfigKeysMatching(toMatch string) (result []string) {
	re := regexp.MustCompile(toMatch)
	for key := range c.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch
func (c *Configuration) RemoveFromPerennialBranches(branchName string) {
	c.SetPerennialBranches(util.RemoveStringFromSlice(c.GetPerennialBranches(), branchName))
}

// RemoveGitAlias removes the given Git alias.
func (c *Configuration) RemoveGitAlias(command string) *command.Result {
	return c.removeGlobalConfigValue("alias." + command)
}

func (c *Configuration) removeGlobalConfigValue(key string) *command.Result {
	delete(c.globalConfigCache, key)
	return command.RunInDir(c.localDir, "git", "config", "--global", "--unset", key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *Configuration) removeLocalConfigValue(key string) {
	command.RunInDir(c.localDir, "git", "config", "--unset", key)
	delete(c.localConfigCache, key)
}

// RemoveLocalGitConfiguration removes all Git Town configuration
func (c *Configuration) RemoveLocalGitConfiguration() {
	command.RunInDir(c.localDir, "git", "config", "--remove-section", "git-town").OutputSanitized()
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration
func (c *Configuration) RemoveOutdatedConfiguration() {
	for child, parent := range c.GetParentBranchMap() {
		if !HasBranch(child) || !HasBranch(parent) {
			c.DeleteParentBranch(child)
		}
	}
}

func (c *Configuration) setGlobalConfigValue(key, value string) *command.Result {
	c.globalConfigCache[key] = value
	return command.RunInDir(c.localDir, "git", "config", "--global", key, value)
}

// setConfigurationValue sets the local configuration with the given key to the given value.
func (c *Configuration) setLocalConfigValue(key, value string) *command.Result {
	c.localConfigCache[key] = value
	return command.RunInDir(c.localDir, "git", "config", key, value)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Configuration) SetMainBranch(branchName string) *command.Result {
	return c.setLocalConfigValue("git-town.main-branch-name", branchName)
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) SetNewBranchPush(value bool, global bool) *command.Result {
	if global {
		return c.setGlobalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
	}
	return c.setLocalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// SetOffline updates whether Git Town is in offline mode
func (c *Configuration) SetOffline(value bool) *command.Result {
	return c.setGlobalConfigValue("git-town.offline", strconv.FormatBool(value))
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Configuration) SetParentBranch(branchName, parentBranchName string) *command.Result {
	return c.setLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func (c *Configuration) SetPerennialBranches(branchNames []string) *command.Result {
	return c.setLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Configuration) SetPullBranchStrategy(strategy string) *command.Result {
	return c.setLocalConfigValue("git-town.pull-branch-strategy", strategy)
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) ShouldNewBranchPush() bool {
	config := c.getLocalOrGlobalConfigValue("git-town.new-branch-push-flag")
	if config == "" {
		return false
	}
	return util.StringToBool(config)
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches up to the origin remote.
func (c *Configuration) ShouldNewBranchPushGlobal() bool {
	config := c.getGlobalConfigValue("git-town.new-branch-push-flag")
	return config == "true"
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (c *Configuration) ShouldSyncUpstream() bool {
	return c.getLocalOrGlobalConfigValue("git-town.sync-upstream") != "false"
}

// ValidateIsOnline asserts that Git Town is not in offline mode
func (c *Configuration) ValidateIsOnline() error {
	if c.IsOffline() {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
