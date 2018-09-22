package git

// IVersionService provides methods about the git version
type IVersionService interface {
	EnsureVersionRequirementSatisfied()
}
