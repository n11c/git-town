package steps

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {

	suite.Step(`^Git Town is no longer configured for this repository$`, func() error {
		output, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "config", "--local", "--get-regex", "git-town")
		exitError := err.(*exec.ExitError)
		if exitError.ExitCode() != 1 {
			return errors.New("git config should return exit code 1 if no matching configuration found")
		}
		if strings.TrimSpace(output) != "" {
			return errors.Wrapf(err, "expeced no local Git Town configuration but got %q", output)
		}
		return nil
	})

	suite.Step(`^I haven\'t configured Git Town yet$`, func() error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().DeleteMainBranchConfiguration()
		if err != nil {
			return errors.Wrap(err, "cannot delete main branch config")
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().DeletePerennialBranchConfiguration()
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		branches := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is set to "(true|false)"$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.Wrapf(err, "cannot parse %q into bool", value)
		}
		outcome := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().SetNewBranchPush(b, false)
		return outcome.Err()
	})

	suite.Step(`^the main branch is configured as "([^"]+)"$`, func(name string) error {
		outcome := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().SetMainBranch(name)
		return outcome.Err()
	})

	suite.Step(`^the main branch is now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetMainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the main branch name is not configured$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().DeleteMainBranchConfiguration()
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().DeletePerennialBranchConfiguration()
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		outcome := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().AddToPerennialBranches(branch1, branch2)
		return outcome.Err()
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expeced 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expeced 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 {
			return fmt.Errorf("expected %q, got %q", branch1, actual)
		}
		if actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", branch2, actual)
		}
		return nil
	})

	suite.Step(`^there are now no perennial branches$`, func() error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) > 0 {
			return fmt.Errorf("expeced no perennial branches, got %q", actual)
		}
		return nil
	})
}
