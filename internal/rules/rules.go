package rules

import (
	"github.com/wmatex/tmux-workspace/internal/projects"
	"github.com/wmatex/tmux-workspace/internal/tmux"
	"github.com/wmatex/tmux-workspace/internal/utils"
)

type Rules struct {
	Overrides map[string]Rule
	Rules     []Rule
}

type Rule struct {
	Checks  []RuleCheck
	Name    string
	Hooks   map[string][]string
	Windows []tmux.Window
}

func (r *Rules) GetSatisfied(p *projects.Project) []*Rule {
	var valid []*Rule

	override, ok := r.Overrides[p.Name]
	if ok {
		valid = append(valid, &override)
		return valid
	}

	for _, rule := range r.Rules {
		satisfies := true

		for _, check := range rule.Checks {
			if !check.IsSatisfiedForProject(p, valid) {
				satisfies = false
				break
			}
		}

		if satisfies {
			valid = append(valid, &rule)
		}
	}

	return valid
}

func MergeWindows(rules []*Rule) []*tmux.Window {
	var windows []*tmux.Window
	for _, r := range rules {
		for _, w := range r.Windows {
			windows = append(windows, &w)
		}
	}

	return utils.Merge(windows)
}
