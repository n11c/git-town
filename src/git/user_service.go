package git

import "github.com/Originate/git-town/src/command"

type userService struct{}

// NewUserService returns a new IUserService
func NewUserService() IUserService {
	return &userService{}
}

// GetLocalAuthor returns the locally Git configured user
func (u *userService) GetLocalAuthor() string {
	name := command.New("git", "config", "user.name").Output()
	email := command.New("git", "config", "user.email").Output()
	return name + " <" + email + ">"
}
