package projects

import (
	"slices"
	"strings"

	"github.com/wmatex/automux/internal/cmd_exec/tmux"
)

type Projects struct {
	Map        map[string]*Project
	sortedKeys []string
}

func (p *Projects) MergeProjectsWithSessions(sessions []tmux.Session) []*Project {
	var result []*Project

	for _, s := range sessions {
		_, ok := (*p).Map[s.Name]
		if !ok {
			p.Map[s.Name] = &Project{
				Name:    s.Name,
				Path:    s.Path,
				Active:  s.Active,
				Running: true,
			}
		} else {
			p.Map[s.Name].Path = s.Path
			p.Map[s.Name].Active = s.Active
			p.Map[s.Name].Running = true
		}
	}

	var names []string
	for n := range p.Map {
		names = append(names, n)
	}
	slices.SortFunc(names, func(a, b string) int {
		pA := p.Map[a]
		pB := p.Map[b]

		if pA.Running && pB.Running {
			return strings.Compare(pA.Name, pB.Name)
		}

		if pA.Running {
			return -1
		}

		if pB.Running {
			return 1
		}

		return strings.Compare(pA.Name, pB.Name)
	})

	for _, name := range names {
		project := p.Map[name]
		result = append(result, project)
	}

	return result
}
