package tmux

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wmatex/automux/internal/cmd_exec"
)

const SEPARATOR = "|"

type Session struct {
	Name         string
	Active       bool
	Path         string
	LastActivity int
}

func createCmdBuilder(args []string) *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("tmux", args)
}

func NewSession(name, path string) (error, int) {
	if path != "" {
		return createCmdBuilder([]string{"new-session", "-d", "-s", name, "-c", path}).Exec()
	} else {
		return createCmdBuilder([]string{"new-session", "-d", "-s", name}).Exec()
	}
}

func SwitchToSession(name string) (error, int) {
	return createCmdBuilder([]string{"switch-client", "-t", name}).Exec()
}

func GetActiveSessions() ([]Session, error) {
	cmd := createCmdBuilder([]string{
		"list-sessions",
		"-F",
		fmt.Sprintf("#S%[1]s#{session_attached}%[1]s#{session_path}%[1]s#{session_activity}", SEPARATOR),
	})

	output, err, _ := cmd.ExecWithOutput()
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

		lastActivity, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Fatalf("cannot parse last activity time: %s\n", err)
		}
		session := Session{
			Name:         parts[0],
			Active:       active,
			Path:         parts[2],
			LastActivity: lastActivity,
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
