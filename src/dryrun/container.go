package dryrun

import (
	"github.com/Originate/exit"
	"go.uber.org/dig"
)

func ProvideServices(container *dig.Container) {
	exit.If(container.Provide(NewStateService))
}
