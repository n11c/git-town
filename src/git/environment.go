package git

import (
	"errors"

	"github.com/Originate/git-town/src/command"
)

// ValidateIsRepository asserts that the current directory is in a repository
func ValidateIsRepository() error {
	if IsRepository() {
		return nil
	}
	return errors.New("this is not a Git repository")
}

// isRepository is cached in order to minimize the number of git commands run
var isRepository bool
var isRepositoryInitialized bool

// IsRepository returns whether or not the current directory is in a repository
func IsRepository() bool {
	if !isRepositoryInitialized {
		isRepository = command.Run("git", "rev-parse").Err() == nil
		isRepositoryInitialized = true
	}
	return isRepository
}
