package rules

import (
	"github.com/wmatex/automux/internal/projects"
	"github.com/wmatex/automux/internal/tmux"
	"github.com/wmatex/automux/internal/utils"
)

type Rules struct {
	Rules     []Rule
	Overrides map[string]Rule
}

type Rule struct {
	Checks  []RuleCheck
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
			if !check.IsSatisfiedForProject(p) {
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
