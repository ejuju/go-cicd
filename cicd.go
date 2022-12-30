package cicd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	logger           io.Writer
	echoExecCommands bool
}

func NewRunner() *Runner {
	return &Runner{
		logger:           os.Stdout,
		echoExecCommands: true,
	}
}

func (r *Runner) Run(steps ...*Step) error {
	for _, step := range steps {
		fmt.Fprintf(r.logger, "Step \"%s\"\n", step.name)
		for i, fn := range step.fns {
			err := fn()
			if err != nil {
				return err
			}
			fmt.Fprintf(r.logger, "Completed %d/%d\n", i+1, len(step.fns))
		}
	}
	fmt.Fprintf(r.logger, "OK\n")
	return nil
}

type Step struct {
	name string
	fns  []func() error
}

func NewStep(name string, fns ...func() error) *Step {
	return &Step{name: name, fns: fns}
}

func Exec(cmdstr string) func() error {
	parts := strings.Split(cmdstr, " ")

	return func() error {
		cmd := exec.Command(parts[0], parts[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(output)
		return err
	}
}

func SetEnv(key string, value string) func() error {
	return func() error { return os.Setenv(key, value) }
}
