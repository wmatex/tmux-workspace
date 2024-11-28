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
				Name:       s.Name,
				Path:       s.Path,
				Active:     s.Active,
				Running:    true,
				LastActive: s.LastActivity,
			}
		} else {
			p.Map[s.Name].Path = s.Path
			p.Map[s.Name].Active = s.Active
			p.Map[s.Name].Running = true
			p.Map[s.Name].LastActive = s.LastActivity
		}
	}

	var names []string
	for n := range p.Map {
		names = append(names, n)
	}
	slices.SortFunc(names, func(a, b string) int {
		pA := p.Map[a]
		pB := p.Map[b]

		if pA.LastActive == pB.LastActive {
			return strings.Compare(pA.Name, pB.Name)
		}

		return pB.LastActive - pA.LastActive
	})

	for _, name := range names {
		project := p.Map[name]
		if !project.Active {
			result = append(result, project)
		}
	}

	return result
}
