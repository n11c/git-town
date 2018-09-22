package git

import (
	"github.com/Originate/exit"
	"go.uber.org/dig"
)

func ProvideServices(container *dig.Container) {
	exit.If(container.Provide(NewBranchService))
	exit.If(container.Provide(NewConfigService))
	exit.If(container.Provide(NewCurrentBranchService))
	exit.If(container.Provide(NewEnvironmentService))
	exit.If(container.Provide(NewLogService))
	exit.If(container.Provide(NewPrintableService))
	exit.If(container.Provide(NewShaService))
	exit.If(container.Provide(NewSquashMergeService))
	exit.If(container.Provide(NewStatusService))
	exit.If(container.Provide(NewUserService))
	exit.If(container.Provide(NewVersionService))
}
