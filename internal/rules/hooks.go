package rules

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/wmatex/tmux-workspace/internal/cmd_exec"
	"github.com/wmatex/tmux-workspace/internal/projects"
	"github.com/wmatex/tmux-workspace/internal/tmux"
	"github.com/wmatex/tmux-workspace/internal/utils"
)

const (
	START uint8 = iota
	END   uint8 = iota
)

type Hook struct {
	Cmd       string
	Lifecycle uint8
}

func (h *Hook) Id() string {
	return fmt.Sprintf("%d:%s", h.Lifecycle, h.Cmd)
}

func (h *Hook) Merge(o *Hook) {
	// No action required
}

func SetupHooks(p *projects.Project, rules []*Rule) error {
	merged := mergeHooks(START, rules)

	if len(merged) > 0 {
		err, _ := tmux.Popup(p.Path, fmt.Sprintf("%s session start -p %s", os.Args[0], p.Name))
		if err != nil {
			return err
		}
	}

	err, _ := tmux.SetHook("session-closed", fmt.Sprintf("%s session end -p #{hook_session_name}", os.Args[0]))
	return err
}

func RunHooks(p *projects.Project, lifecycle uint8, rules []*Rule) error {
	merged := mergeHooks(lifecycle, rules)
	ch := make(chan int, len(merged))

	for _, hook := range merged {
		go runHook(p.Path, hook, ch)
	}

	ok := true
	for _, hook := range merged {
		status := <-ch

		if status != 0 {
			ok = false
			log.Printf("hook %s failed\n", hook.Cmd)
		}
	}

	if !ok {
		return errors.New("start hooks failed")
	}

	return nil
}

func mergeHooks(lifecycle uint8, rules []*Rule) []*Hook {
	var filtered []*Hook

	var key string
	switch lifecycle {
	case START:
		key = "start"

	case END:
		key = "end"
	}

	for _, r := range rules {
		for lc, cmds := range r.Hooks {
			if lc == key {
				for _, cmd := range cmds {
					filtered = append(filtered, &Hook{
						Lifecycle: lifecycle,
						Cmd:       cmd,
					})
				}
			}
		}
	}

	return utils.Merge(filtered)
}

func runHook(dir string, hook *Hook, ch chan int) {
	cmdArgs := utils.SplitArgs(hook.Cmd)
	_, status := cmd_exec.
		NewCmdExec(cmdArgs[0], cmdArgs[1:]).
		SetWorkingDirectory(dir).
		Exec(true)

	ch <- status
}
