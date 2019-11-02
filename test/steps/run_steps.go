package steps

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/test"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
func RunSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I run "([^"]+)"$`, func(command string) error {
		fs.activeScenarioState.lastRunOutput, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunString(command)
		return nil
	})

	suite.Step(`^I run "([^"]+)" and answer the prompts:$`, func(cmd string, input *gherkin.DataTable) error {
		fs.activeScenarioState.lastRunOutput, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.DeveloperRepo.RunStringWithInput(cmd, tableToInput(input))
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *gherkin.DataTable) error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunOutput)
		table := test.RenderExecutedGitCommands(commands, input)
		diff, errorCount := table.Equal(input)
		if errorCount != 0 {
			return fmt.Errorf("found %d differences:\n%s", errorCount, diff)
		}
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := test.GitCommandsInGitTownOutput(fs.activeScenarioState.lastRunOutput)
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})
}

func tableToInput(table *gherkin.DataTable) []command.Input {
	var result []command.Input
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answer := row.Cells[1].Value
		answer = strings.ReplaceAll(answer, "[ENTER]", "\n")
		answer = strings.ReplaceAll(answer, "[DOWN]", "\x1b[B")
		answer = strings.ReplaceAll(answer, "[UP]", "\x1b[A")
		answer = strings.ReplaceAll(answer, "[SPACE]", " ")
		result = append(result, command.Input{Prompt: row.Cells[0].Value, Answer: answer})
	}
	return result
}
