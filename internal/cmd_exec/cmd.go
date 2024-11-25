package cmd_exec

import (
	"os/exec"
	"strings"
)

type CmdExecBuilder struct {
	prog  string
	args  []string
	input []string
}

func NewCmdExec(prog string, args []string) *CmdExecBuilder {
	c := &CmdExecBuilder{prog, args, []string{}}

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

func (c *CmdExecBuilder) exec() *exec.Cmd {
	cmd := exec.Command(c.prog, c.args...)
	if len(c.input) > 0 {
		cmd.Stdin = strings.NewReader(strings.Join(c.input, "\n"))
	}

	return cmd
}

func (c *CmdExecBuilder) ExecWithOutput() (string, error) {
	cmd := c.exec()

	out, err := cmd.Output()
	return string(out), err
}

func (c *CmdExecBuilder) Exec() error {
	cmd := c.exec()

	return cmd.Run()
}
