package browsers

// IOpenService represents a class that helps determine what browser to use to open a url
type IOpenService interface {
	GetOpenBrowserCommand() string
}
