package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
)

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I have Git "([^"]*)" installed$`, func(version string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.AddTempShellOverride(
			"git",
			fmt.Sprintf(`#!/usr/bin/env bash
		if [ "$1" = "version" ]; then
			echo "git version %s"
		fi`, version))
		return err
	})
}
