package git

import (
	"errors"

	"github.com/Originate/git-town/src/command"
)

// EnvironmentService provides methods to query about the git environment
type environmentService struct {
	configService           IConfigService
	isRepository            bool // cached to minimize the number of git commands run
	isRepositoryInitialized bool
	rootDirectory           string // cached to minimize the number of git commands run
}

// NewEnvironmentService returns a new IEnvironmentService
func NewEnvironmentService(configService IConfigService) IEnvironmentService {
	return &environmentService{
		configService: configService,
	}
}

// GetRootDirectory returns the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (e *environmentService) GetRootDirectory() string {
	if e.rootDirectory == "" {
		e.rootDirectory = command.New("git", "rev-parse", "--show-toplevel").Output()
	}
	return e.rootDirectory
}

// IsRepository returns whether or not the current directory is in a repository
func (e *environmentService) IsRepository() bool {
	if !e.isRepositoryInitialized {
		e.isRepository = command.New("git", "rev-parse").Err() == nil
		e.isRepositoryInitialized = true
	}
	return e.isRepository
}

// ValidateIsRepository asserts that the current directory is in a repository
func (e *environmentService) ValidateIsRepository() error {
	if e.IsRepository() {
		return nil
	}
	return errors.New("This is not a Git repository")
}
