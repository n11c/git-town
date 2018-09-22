package git

// IEnvironmentService provides methods to query about the git environment
type IEnvironmentService interface {
	GetRootDirectory() string
	IsRepository() bool
	ValidateIsRepository() error
}
