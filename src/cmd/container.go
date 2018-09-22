package cmd

import (
	"github.com/Originate/git-town/src/browsers"
	"github.com/Originate/git-town/src/dryrun"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/steps/listbuilder"
	"go.uber.org/dig"
)

func GetContainer() *dig.Container {
	container := dig.New()
	browsers.ProvideServices(container)
	dryrun.ProvideServices(container)
	git.ProvideServices(container)
	listbuilder.ProvideServices(container)
	prompt.ProvideServices(container)
	script.ProvideServices(container)
	steps.ProvideServices(container)
	return container
}
