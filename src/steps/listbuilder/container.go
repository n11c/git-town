package listbuilder

import (
	"github.com/Originate/exit"
	"go.uber.org/dig"
)

func ProvideServices(container *dig.Container) {
	exit.If(container.Provide(NewAppendStepListBuilderService))
	exit.If(container.Provide(NewSyncStepListBuilderService))
}
