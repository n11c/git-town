package git

// ILogService provides methods for the git log
type ILogService interface {
	GetLastCommitMessage() string
}
