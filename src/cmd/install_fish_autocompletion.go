package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Originate/git-town/src/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var installFishAutocompletionCommand = &cobra.Command{
	Use:   "install-fish-autocompletion",
	Short: "Installs the autocompletion definition for Fish shell",
	Run: func(cmd *cobra.Command, args []string) {
		err := installFishAutocompletion()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
}

func installFishAutocompletion() error {
	folderName := path.Join(os.Getenv("HOME"), ".config", "fish", "completions")
	filename := path.Join(folderName, "git.fish")
	err := os.MkdirAll(path.Dir(filename), 0700)
	if err != nil {
		return errors.Wrapf(err, "cannot create folder %q", folderName)
	}
	if util.DoesFileExist(filename) {
		util.ExitWithErrorMessage("Git autocompletion for Fish shell already exists")
	}
	err = ioutil.WriteFile(filename, []byte(buildAutocompletionDefinition()), 0644)
	if err != nil {
		return errors.Wrapf(err, "cannot write file %q", filename)
	}
	fmt.Println("Git autocompletion for Fish shell installed")
	return nil
}

var fishAutocompletionTemplate = `
# All Git Town commands
set git_town_commands %s

# Indicates through its error code whether the command line to auto-complete
# already contains a Git Town command or not.
#
# - doesn't have command yet: exit code 0
# - has command already: exit code 1
function __fish_complete_git_town_no_command
  for cmd in (commandline -otc)
    if contains $cmd $git_town_commands
      return 1
    end
  end
  return 0
end


# Define autocompletion for the Git Town commands themselves.
#
# These only get autocompleted if there is no Git Town command present in the
# command line already.
# This is done through __fish_complete_git_town_no_command
%s

# Define autocompletion of Git branch names.
#
# This is only enabled for commands that take branch names.
# This is achieved through __fish_complete_git_town_command_takes_branch
complete --command git --arguments "(git branch | tr -d '* ')" --no-files
`

type autocompleteDefinition struct {
	name        string
	description string
}

func buildAutocompletionDefinition() string {
	commands := []autocompleteDefinition{
		{name: "abort", description: abortCmd.Short},
		{name: "continue", description: configCommand.Short},
		{name: "hack", description: hackCmd.Short},
		{name: "kill", description: killCommand.Short},
		{name: "new-pull-request", description: newPullRequestCommand.Short},
		{name: "prune-branches", description: pruneBranchesCommand.Short},
		{name: "rename-branch", description: renameBranchCommand.Short},
		{name: "repo", description: repoCommand.Short},
		{name: "ship", description: shipCmd.Short},
		{name: "sync", description: syncCmd.Short},
		{name: "undo", description: undoCmd.Short},
	}

	commandsSpaceSeparated := ""
	for _, command := range commands {
		commandsSpaceSeparated += command.name + " "
	}
	commandAutocompletion := ""
	for _, command := range commands {
		commandAutocompletion += fmt.Sprintf("complete --command git --arguments '%s' --description '%s' --condition '__fish_complete_git_town_no_command' --no-files\n", command.name, command.description)
	}

	return fmt.Sprintf(fishAutocompletionTemplate, commandsSpaceSeparated, commandAutocompletion)
}

func init() {
	RootCmd.AddCommand(installFishAutocompletionCommand)
}
