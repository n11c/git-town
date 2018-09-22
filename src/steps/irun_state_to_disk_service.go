package steps

// IRunStateToDiskService represents a type that deletes, loads, and saves run state to disk
type IRunStateToDiskService interface {
	LoadPreviousRunState() *RunState
	DeletePreviousRunState()
	SaveRunState(runState *RunState)
}
