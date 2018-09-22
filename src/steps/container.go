package steps

import (
	"github.com/Originate/exit"
	"go.uber.org/dig"
)

func ProvideServices(container *dig.Container) {
	exit.If(container.Provide(NewRunService))
	exit.If(container.Provide(NewRunStateToDiskService))
}
