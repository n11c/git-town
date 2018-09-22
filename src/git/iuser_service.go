package git

// IUserService provides methods about the git user
type IUserService interface {
	GetLocalAuthor() string
}
