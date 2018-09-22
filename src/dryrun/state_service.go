package dryrun

// StateService holds the state of dry run
type StateService struct {
	currentBranchName string
	isActive          bool
}

// NewStateService returns a new IStateService
func NewStateService() IStateService {
	return &StateService{}
}

// Activate enables dry-run mode
func (s *StateService) Activate(initialBranchName string) {
	s.isActive = true
	s.SetCurrentBranchName(initialBranchName)
}

// IsActive returns whether of not dry-run mode is active
func (s *StateService) IsActive() bool {
	return s.isActive
}

// GetCurrentBranchName returns the current branch name for dry-run mode
func (s *StateService) GetCurrentBranchName() string {
	return s.currentBranchName
}

// SetCurrentBranchName sets the current branch name for dry-run mode
func (s *StateService) SetCurrentBranchName(branchName string) {
	s.currentBranchName = branchName
}
