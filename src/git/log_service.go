package git

import "github.com/Originate/git-town/src/command"

// LogService provides methods for the git log
type LogService struct{}

func NewLogService() ILogService {
	return &LogService{}
}

// GetLastCommitMessage returns the commit message for the last commit
func (l *LogService) GetLastCommitMessage() string {
	return command.New("git", "log", "-1", "--format=%B").Output()
}
