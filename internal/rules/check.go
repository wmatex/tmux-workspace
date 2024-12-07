package rules

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/wmatex/automux/internal/cmd_exec"
	"github.com/wmatex/automux/internal/projects"
	"github.com/wmatex/automux/internal/utils"
)

type RuleCheck interface {
	IsSatisfiedForProject(p *projects.Project) bool
}

type DirExistsRule struct {
	dirPath string
}

func (r *DirExistsRule) IsSatisfiedForProject(p *projects.Project) bool {
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

func (r *FileExistsRule) IsSatisfiedForProject(p *projects.Project) bool {
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

func (r *ExecRule) IsSatisfiedForProject(p *projects.Project) bool {
	args := utils.SplitArgs(r.cmd)

	_, status := cmd_exec.
		NewCmdExec(args[0], args[1:]).
		SetWorkingDirectory(p.Path).
		Exec(false)

	return status == 0
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
	}

	return nil, errors.New(fmt.Sprintf("undefined rule '%s'", ruleName))
}
