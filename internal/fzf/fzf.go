package fzf

import (
	"strings"

	"github.com/wmatex/tmux-workspace/internal/cmd_exec"
	"github.com/wmatex/tmux-workspace/internal/projects"
)

const (
	DEFAULT_ACTION  = iota
	SKIP_START_HOOK = iota
)

type PickResult struct {
	Query       string
	ProjectName string
	Action      int
}

func createCmdBuilder() *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("fzf", []string{
		"--no-sort",
		"--print-query",
		"--tmux", "center",
		"--bind", "enter:accept-or-print-query",
		"--bind", "ctrl-o:print(SKIP_START_HOOK)+accept-or-print-query",
		"--bind", "esc:abort",
	})
}

func parseProjectName(line string) string {
	parts := strings.Split(line, " ")
	return parts[1]
}

func ProjectPick(projects []*projects.Project) (*PickResult, error) {
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
			return nil, nil
		} else if code != 1 {
			return nil, err
		}
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	pick := PickResult{}
	pick.Query = lines[0]

	if len(lines) == 1 {
		pick.ProjectName = lines[0]
		pick.Action = DEFAULT_ACTION
	} else if len(lines) == 2 {
		pick.ProjectName = parseProjectName(lines[1])
		pick.Action = DEFAULT_ACTION
	} else {
		pick.ProjectName = parseProjectName(lines[2])
		switch lines[1] {
		case "SKIP_START_HOOK":
			pick.Action = SKIP_START_HOOK
		default:
			pick.Action = DEFAULT_ACTION
		}
	}

	return &pick, nil
}
