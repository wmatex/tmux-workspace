package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	args := SplitArgs("ls -l -t /")

	assert.Len(t, args, 4)
	assert.Equal(t, []string{"ls", "-l", "-t", "/"}, args)
}

func TestMultipleSpaces(t *testing.T) {
	args := SplitArgs("ls   -l     -t /")

	assert.Len(t, args, 4)
	assert.Equal(t, []string{"ls", "-l", "-t", "/"}, args)
}

func TestSingleQuotes(t *testing.T) {
	args := SplitArgs("echo 'Hello, World!'")

	assert.Len(t, args, 2)
	assert.Equal(t, []string{"echo", "Hello, World!"}, args)
}

func TestDoubleQuotes(t *testing.T) {
	args := SplitArgs("echo \"Hello, World!\"")

	assert.Len(t, args, 2)
	assert.Equal(t, []string{"echo", "Hello, World!"}, args)
}

func TestMix(t *testing.T) {
	args := SplitArgs("echo \"Hello, World!\"    'Hello,    \"World\"'     -e")

	assert.Len(t, args, 4)
	assert.Equal(t, []string{"echo", "Hello, World!", "Hello,    \"World\"", "-e"}, args)
}
