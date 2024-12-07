package fzf

import (
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
		"--bind", "esc:abort",
	})
}

func ProjectPick(projects []*projects.Project) (string, error) {
	var input []string
	for _, p := range projects {
		input = append(input, (*p).Format())
	}

	output, err, code := createCmdBuilder().
		SetInput(input).
		AddArguments([]string{"-d", "\\s", "--nth", "2", "--ansi", "--highlight-line"}).
		CaptureOutput()

	if err != nil {
		if code == 130 {
			return "", nil
		} else if code != 1 {
			return "", err
		}
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	lastLine := lines[len(lines)-1]

	parts := strings.Split(lastLine, " ")
	if len(parts) < 2 {
		return parts[0], nil
	}

	return parts[1], nil
}
