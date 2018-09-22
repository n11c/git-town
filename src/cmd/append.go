package cmd

import (
	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/steps/listbuilder"
	"github.com/Originate/git-town/src/util"
	"github.com/davecgh/go-spew/spew"

	"github.com/spf13/cobra"
)

type appendConfig struct {
	ParentBranch string
	TargetBranch string
}

var appendCommand = &cobra.Command{
	Use:   "append <branch>",
	Short: "Creates a new feature branch as a child of the current branch",
	Long: `Creates a new feature branch as a direct child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository
if and only if new-branch-push-flag is true,
and brings over all uncommitted changes to the new feature branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		container := GetContainer()
		exit.If(container.Invoke(func(appendStepListBuilderService listbuilder.IAppendStepListBuilderService, runService steps.IRunService) {
			stepList := appendStepListBuilderService.GetAppendStepList(args)
			spew.Dump(stepList)
			runState := steps.NewRunState("append", stepList)
			runService.Run(runState)
		}))
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(appendCommand)
}
