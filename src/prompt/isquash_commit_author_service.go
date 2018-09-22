package prompt

type ISquashCommitAuthorService interface {
	GetSquashCommitAuthor(branchName string) string
}
