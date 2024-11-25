package fzf

import (
	"strings"

	"github.com/wmatex/automux/internal/cmd_exec"
)

func createCmdBuilder() *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("fzf", []string{
		"--no-sort",
		"--print-query",
		"--tmux", "center",
		"--bind", "enter:replace-query+print-query",
	})
}

func ProjectPick(projects []string) (string, error) {
	output, err := createCmdBuilder().
		SetInput(projects).
		ExecWithOutput()

	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	return lines[len(lines)-1], nil
}
