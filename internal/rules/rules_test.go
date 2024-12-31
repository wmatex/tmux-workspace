package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmatex/tmux-workspace/internal/tmux"
)

func TestMergeTwoIdentical(t *testing.T) {
	r := []*Rule{
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
			},
		},
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
			},
		},
	}

	w := MergeWindows(r)

	assert.Len(t, w, 1)
	assert.Equal(t, "editor", w[0].Name)
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, "echo test", w[0].Panes[0].Cmd)
}

func TestMerge3_2identical(t *testing.T) {
	r := []*Rule{
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
			},
		},
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
				{
					Name: "logs",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo logs",
						},
					},
				},
			},
		},
	}

	w := MergeWindows(r)

	assert.Len(t, w, 2)
	assert.Equal(t, "editor", w[0].Name)
	assert.Equal(t, "logs", w[1].Name)
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, "echo test", w[0].Panes[0].Cmd)
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, "echo logs", w[1].Panes[0].Cmd)
}

func TestMerge3_2identical_different_order(t *testing.T) {
	r := []*Rule{
		{
			Windows: []tmux.Window{
				{
					Name: "logs",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo logs",
						},
					},
				},
			},
		},
		{
			Windows: []tmux.Window{
				{
					Name: "logs",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo logs",
						},
					},
				},
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
			},
		},
	}

	w := MergeWindows(r)

	assert.Len(t, w, 2)
	assert.Equal(t, "logs", w[0].Name)
	assert.Equal(t, "editor", w[1].Name)
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, "echo logs", w[0].Panes[0].Cmd)
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, "echo test", w[1].Panes[0].Cmd)
}

func TestMerge3_MergePanes(t *testing.T) {
	r := []*Rule{
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo test",
						},
					},
				},
			},
		},
		{
			Windows: []tmux.Window{
				{
					Name: "editor",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo another test",
						},
					},
				},
				{
					Name: "logs",
					Panes: []*tmux.Pane{
						{
							Cmd: "echo logs",
						},
					},
				},
			},
		},
	}

	w := MergeWindows(r)

	assert.Len(t, w, 2)
	assert.Equal(t, "editor", w[0].Name)
	assert.Equal(t, "logs", w[1].Name)
	assert.Len(t, w[0].Panes, 2)
	assert.Equal(t, "echo test", w[0].Panes[0].Cmd)
	assert.Equal(t, "echo another test", w[0].Panes[1].Cmd)
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, "echo logs", w[1].Panes[0].Cmd)
}
