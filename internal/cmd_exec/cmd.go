package cmd_exec

import (
	"os"
	"os/exec"
	"strings"
)

type CmdExecBuilder struct {
	prog  string
	args  []string
	input []string
	dir   string
}

func NewCmdExec(prog string, args []string) *CmdExecBuilder {
	c := &CmdExecBuilder{prog, args, []string{}, ""}

	return c
}

func (c *CmdExecBuilder) AddArguments(args []string) *CmdExecBuilder {
	c.args = append(c.args, args...)
	return c
}

func (c *CmdExecBuilder) SetInput(lines []string) *CmdExecBuilder {
	c.input = lines

	return c
}

func (c *CmdExecBuilder) SetWorkingDirectory(dir string) *CmdExecBuilder {
	c.dir = dir

	return c
}

func (c *CmdExecBuilder) exec() *exec.Cmd {
	cmd := exec.Command(c.prog, c.args...)
	if c.dir != "" {
		cmd.Dir = c.dir
	}
	if len(c.input) > 0 {
		cmd.Stdin = strings.NewReader(strings.Join(c.input, "\n"))
	}
	cmd.Stderr = os.Stderr

	return cmd
}

func (c *CmdExecBuilder) CaptureOutput() (string, error, int) {
	cmd := c.exec()

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr, exitErr.ExitCode()
		} else {
			return string(out), err, 1
		}
	}

	return string(out), err, 0
}

func (c *CmdExecBuilder) Exec(passOutput bool) (error, int) {
	cmd := c.exec()
	if passOutput {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr, exitErr.ExitCode()
		} else {
			return err, 1
		}
	}

	return nil, 0
}
