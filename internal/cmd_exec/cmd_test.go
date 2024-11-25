package cmd_exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	prog := "echo"
	args := []string{"hello"}
	cmd := NewCmdExec(prog, args)

	err := cmd.Exec()
	assert.Nil(t, err)
}

func TestExecWithOutput(t *testing.T) {
	prog := "echo"
	args := []string{"hello"}
	cmd := NewCmdExec(prog, args)

	output, err := cmd.ExecWithOutput()
	assert.Nil(t, err)

	expectedOutput := "hello\n"
	assert.Equalf(t, expectedOutput, output, "expected %s, got %s", expectedOutput, output)
}
