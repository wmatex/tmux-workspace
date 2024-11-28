package cmd_exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	prog := "echo"
	args := []string{"hello"}
	cmd := NewCmdExec(prog, args)

	err, code := cmd.Exec()
	assert.Nil(t, err)
	assert.Equal(t, code, 0)
}

func TestExecWithOutput(t *testing.T) {
	prog := "echo"
	args := []string{"hello"}
	cmd := NewCmdExec(prog, args)

	output, err, code := cmd.ExecWithOutput()
	assert.Nil(t, err)
	assert.Equal(t, code, 0)

	expectedOutput := "hello\n"
	assert.Equalf(t, expectedOutput, output, "expected %s, got %s", expectedOutput, output)
}
