package command

import (
	"io"
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

// Input contains the user input for a subshell command.
type Input struct {
	Prompt string
	Answer string
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
	result := Result{
		cmd:  opts.Cmd,
		args: opts.Args,
	}
	if len(opts.Input) == 0 {
		output, err := subProcess.CombinedOutput()
		result.err = err
		result.output = string(output)
		return &result
	}

	// here we have to run with opts.Input set
	var input io.WriteCloser
	input, result.err = subProcess.StdinPipe()
	if result.err != nil {
		return &result
	}
	var output io.ReadCloser
	output, result.err = subProcess.StdoutPipe()
	if result.err != nil {
		return &result
	}
	scanner := NewByteStreamScanner(output)
	result.err = subProcess.Start()
	if result.err != nil {
		return &result
	}
	for i := range opts.Input {
		<-scanner.WaitForText(opts.Input[i].Prompt)
		input.Write([]byte(opts.Input[i].Answer))
	}
	result.err = subProcess.Wait()
	result.output = scanner.ReceivedText()
	return &result
}
