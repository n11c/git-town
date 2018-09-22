package steps

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// RunStateToDiskService deletes, loads, and saves run state to disk
type RunStateToDiskService struct {
	gitEnvironmentService git.IEnvironmentService
}

// NewRunStateToDiskService returns a new IRunStateToDiskService
func NewRunStateToDiskService(gitEnvironmentService git.IEnvironmentService) IRunStateToDiskService {
	return &RunStateToDiskService{
		gitEnvironmentService: gitEnvironmentService,
	}
}

// LoadPreviousRunState loads the run state from disk if it exists or creates a new run state
func (r *RunStateToDiskService) LoadPreviousRunState() *RunState {
	filename := r.getRunResultFilename()
	if util.DoesFileExist(filename) {
		var runState RunState
		content, err := ioutil.ReadFile(filename)
		exit.If(err)
		err = json.Unmarshal(content, &runState)
		exit.If(err)
		return &runState
	}
	return nil
}

// DeletePreviousRunState deletes the previous run state from disk
func (r *RunStateToDiskService) DeletePreviousRunState() {
	filename := r.getRunResultFilename()
	if util.DoesFileExist(filename) {
		exit.If(os.Remove(filename))
	}
}

// SaveRunState saves the run state to disk
func (r *RunStateToDiskService) SaveRunState(runState *RunState) {
	content, err := json.MarshalIndent(runState, "", "  ")
	exit.If(err)
	filename := r.getRunResultFilename()
	err = ioutil.WriteFile(filename, content, 0644)
	exit.If(err)
}

func (r *RunStateToDiskService) getRunResultFilename() string {
	replaceCharacterRegexp, err := regexp.Compile("[[:^alnum:]]")
	exit.IfWrap(err, "Error compiling replace character expression")
	directory := replaceCharacterRegexp.ReplaceAllString(r.gitEnvironmentService.GetRootDirectory(), "-")
	return path.Join(os.TempDir(), directory)
}
