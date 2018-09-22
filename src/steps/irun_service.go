package steps

type IRunService interface {
	Run(runState *RunState)
}
