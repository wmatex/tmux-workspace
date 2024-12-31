package rules

import (
	"fmt"
	"os"
	"path"

	"github.com/wmatex/tmux-workspace/internal/cmd_exec"
	"github.com/wmatex/tmux-workspace/internal/projects"
	"github.com/wmatex/tmux-workspace/internal/utils"
)

type RuleCheck interface {
	IsSatisfiedForProject(p *projects.Project, valid []*Rule) bool
}

type DirExistsRule struct {
	dirPath string
}

func (r *DirExistsRule) IsSatisfiedForProject(p *projects.Project, valid []*Rule) bool {
	dirPath := path.Join(p.Path, r.dirPath)
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

type FileExistsRule struct {
	filePath string
}

func (r *FileExistsRule) IsSatisfiedForProject(p *projects.Project, valid []*Rule) bool {
	filePath := path.Join(p.Path, r.filePath)
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type ExecRule struct {
	cmd string
}

func (r *ExecRule) IsSatisfiedForProject(p *projects.Project, valid []*Rule) bool {
	args := utils.SplitArgs(r.cmd)

	_, status := cmd_exec.
		NewCmdExec(args[0], args[1:]).
		SetWorkingDirectory(p.Path).
		Exec(false)

	return status == 0
}

type NotActiveRule struct {
	checkName string
}

func (r *NotActiveRule) IsSatisfiedForProject(p *projects.Project, valid []*Rule) bool {
	for _, v := range valid {
		if v.Name == r.checkName {
			return false
		}
	}
	return true
}

func ruleCheckFactory(ruleName, value string) (RuleCheck, error) {
	switch ruleName {
	case "dir_exists":
		return &DirExistsRule{
			dirPath: value,
		}, nil

	case "file_exists":
		return &FileExistsRule{
			filePath: value,
		}, nil
	case "exec":
		return &ExecRule{
			cmd: value,
		}, nil
	case "not_active":
		return &NotActiveRule{
			checkName: value,
		}, nil
	}

	return nil, fmt.Errorf("undefined rule '%s'", ruleName)
}
