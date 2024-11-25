package tmux

import (
	"github.com/wmatex/automux/internal/cmd_exec"
)

func createCmdBuilder(args []string) *cmd_exec.CmdExecBuilder {
	return cmd_exec.NewCmdExec("tmux", args)
}

func NewSession(name string) error {
	return createCmdBuilder([]string{"new-session", "-d", "-s", name}).Exec()
}

func SwitchToSession(name string) error {
	return createCmdBuilder([]string{"switch-client", "-t", name}).Exec()
}
