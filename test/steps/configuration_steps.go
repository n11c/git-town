package steps

import (
	"fmt"
	"strconv"

	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {
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
}
