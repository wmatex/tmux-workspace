package tmux

import (
	"fmt"
	"strings"

	"github.com/wmatex/automux/internal/cmd_exec"
)

const SEPARATOR = "|"

type Session struct {
	Name   string
	Active bool
	Path   string
}

func createCmdBuilder(args []string) *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("tmux", args)
}

func NewSession(name, path string) error {
	return createCmdBuilder([]string{"new-session", "-d", "-s", name, "-c", path}).Exec()
}

func SwitchToSession(name string) error {
	return createCmdBuilder([]string{"switch-client", "-t", name}).Exec()
}

func GetActiveSessions() ([]Session, error) {
	cmd := createCmdBuilder([]string{"list-sessions", "-F", fmt.Sprintf("#S%[1]s#{session_attached}%[1]s#{session_path}", SEPARATOR)})

	output, err := cmd.ExecWithOutput()
	if err != nil {
		return nil, err
	}

	sessions := []Session{}
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		parts := strings.Split(line, SEPARATOR)
		active := false
		if parts[1] == "1" {
			active = true
		}

		session := Session{
			Name:   parts[0],
			Active: active,
			Path:   parts[2],
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
