package command

import (
	"os/exec"
)

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {
	Cmd   string   // the command to run
	Args  []string // arguments for the command
	Dir   string   // the directory in which to execute the command
	Env   []string // environment variables to use, in the format provided by os.Environ()
	Input []Input  // user input to pipe into the command
}

// Run executes the command given in argv notation.
func Run(cmd string, args ...string) *Result {
	return RunWith(Options{Cmd: cmd, Args: args})
}

// RunWith runs the command with the given RunOptions.
func RunWith(opts Options) *Result {
	logRun(opts.Cmd, opts.Args...)
	subProcess := exec.Command(opts.Cmd, opts.Args...) // #nosec
	if opts.Dir != "" {
		subProcess.Dir = opts.Dir
	}
	if opts.Env != nil {
		subProcess.Env = opts.Env
	}
	output, err := subProcess.CombinedOutput()
	return &Result{
		cmd:    opts.Cmd,
		args:   opts.Args,
		err:    err,
		output: string(output),
	}
}
