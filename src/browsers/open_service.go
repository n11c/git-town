package browsers

import (
	"runtime"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// OpenService helps determine what browser to use to open a url
type OpenService struct {
	missingOpenBrowserCommandMessages []string
	openBrowserCommands               []string
}

// NewOpenService returns a new OpenService
func NewOpenService() IOpenService {
	return &OpenService{
		missingOpenBrowserCommandMessages: []string{
			"Cannot open a browser.",
			"If you think this is a bug,",
			"please open an issue at https://github.com/Originate/git-town/issues",
			"and mention your OS and browser.",
		},
		openBrowserCommands: []string{
			"xdg-open",
			"open",
			"cygstart",
			"x-www-browser",
			"firefox",
			"opera",
			"mozilla",
			"netscape",
		},
	}
}

// GetOpenBrowserCommand returns the command to run on the console
// to open the default browser.
func (o *OpenService) GetOpenBrowserCommand() string {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters
		//       like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return "start"
	}
	for _, browserCommand := range openBrowserCommands {
		cmd := command.New("which", browserCommand)
		if cmd.Err() == nil && cmd.Output() != "" {
			return browserCommand
		}
	}
	util.ExitWithErrorMessage(missingOpenBrowserCommandMessages...)
	return ""
}
