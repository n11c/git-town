package git

// IConfigService provides methods around storing configuration settings for the repository.
type IConfigService interface {
	AddToPerennialBranches(branchName string)
	DeleteParentBranch(branchName string)
	EnsureIsFeatureBranch(branchName, errorSuffix string)
	GetAncestorBranches(branchName string) []string
	GetParentBranchMap() map[string]string
	GetChildBranches(branchName string) []string
	GetConfigurationValue(key string) string
	GetGlobalConfigurationValue(key string) string
	GetMainBranch() string
	GetParentBranch(branchName string) string
	GetPerennialBranches() []string
	GetPullBranchStrategy() string
	GetRemoteOriginURL() string
	GetRemoteUpstreamURL() string
	GetURLHostname(url string) string
	GetURLRepositoryName(url string) string
	HasGlobalConfigurationValue(key string) bool
	HasParentBranch(branchName string) bool
	IsAncestorBranch(branchName, ancestorBranchName string) bool
	HasRemote(name string) bool
	IsFeatureBranch(branchName string) bool
	IsMainBranch(branchName string) bool
	IsOffline() bool
	IsPerennialBranch(branchName string) bool
	RemoveAllConfiguration()
	RemoveFromPerennialBranches(branchName string)
	SetMainBranch(branchName string)
	SetParentBranch(branchName, parentBranchName string)
	SetPerennialBranches(branchNames []string)
	SetPullBranchStrategy(strategy string)
	ShouldNewBranchPush() bool
	GetGlobalNewBranchPushFlag() string
	UpdateOffline(value bool)
	UpdateShouldNewBranchPush(value bool)
	UpdateGlobalShouldNewBranchPush(value bool)
	ValidateIsOnline() error
}
