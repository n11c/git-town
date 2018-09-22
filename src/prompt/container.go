package prompt

import (
	"github.com/Originate/exit"
	"go.uber.org/dig"
)

func ProvideServices(container *dig.Container) {
	exit.If(container.Provide(NewBranchService))
	exit.If(container.Provide(NewParentBranchService))
	exit.If(container.Provide(NewSquashCommitAuthorService))
}
