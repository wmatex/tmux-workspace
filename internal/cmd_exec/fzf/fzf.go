package fzf

import (
	"fmt"
	"strings"

	"github.com/wmatex/automux/internal/cmd_exec"
	"github.com/wmatex/automux/internal/projects"
)

func createCmdBuilder() *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("fzf", []string{
		"--no-sort",
		"--print-query",
		"--tmux", "center",
		"--bind", "enter:replace-query+print-query",
	})
}

func ProjectPick(projects []*projects.Project) (string, error) {
	var input []string
	for _, p := range projects {
		input = append(input, fmt.Sprintf("%s %s", (*p).Name, (*p).Path))
	}

	output, err := createCmdBuilder().
		SetInput(input).
		AddArguments([]string{"-d", "\\s", "--nth", "1"}).
		ExecWithOutput()

	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	lastLine := lines[len(lines)-1]

	parts := strings.Split(lastLine, " ")
	return parts[0], nil
}
