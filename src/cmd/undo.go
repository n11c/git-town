package cmd

import (
	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the last run git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		container := GetContainer()
		exit.If(container.Invoke(func(runStateToDiskService steps.IRunStateToDiskService, runService steps.IRunService) {
			runState := runStateToDiskService.LoadPreviousRunState()
			if runState == nil || runState.IsUnfinished() {
				util.ExitWithErrorMessage("Nothing to undo")
			}
			undoRunState := runState.CreateUndoRunState()
			runService.Run(&undoRunState)
		}))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(undoCmd)
}
