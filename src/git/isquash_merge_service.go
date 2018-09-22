package git

// ISquashMergeService provides methods around performing a git squash mere
type ISquashMergeService interface {
	CommentOutSquashCommitMessage(prefix string)
}
