package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmatex/tmux-workspace/internal/projects"
)

func TestDirExistsRule(t *testing.T) {
	p := projects.Project{
		Path: ".",
	}
	rule, err := ruleCheckFactory("dir_exists", "../rules/")
	assert.Nil(t, err)

	result := rule.IsSatisfiedForProject(&p, []*Rule{})
	assert.True(t, result)
}

func TestFileExistsRule(t *testing.T) {
	p := projects.Project{
		Path: ".",
	}
	rule, err := ruleCheckFactory("file_exists", "./check_test.go")
	assert.Nil(t, err)

	result := rule.IsSatisfiedForProject(&p, []*Rule{})
	assert.True(t, result)
}

func TestExecRule(t *testing.T) {
	p := projects.Project{
		Path: ".",
	}
	rule, err := ruleCheckFactory("exec", "jq -e '.array[] | .a' ../../test/sample.json")
	assert.Nil(t, err)

	result := rule.IsSatisfiedForProject(&p, []*Rule{})
	assert.True(t, result)
}

func TestNotActiveRule(t *testing.T) {
	p := projects.Project{
		Path: ".",
	}

	active := []*Rule{
		{Name: "active"},
	}
	rule, err := ruleCheckFactory("not_active", "active")
	assert.Nil(t, err)

	result := rule.IsSatisfiedForProject(&p, active)
	assert.False(t, result)
}

func TestNonExistentRule(t *testing.T) {
	rule, err := ruleCheckFactory("nil", "nil")
	assert.NotNil(t, err)
	assert.Nil(t, rule)
}
