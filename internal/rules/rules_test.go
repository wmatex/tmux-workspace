package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmatex/automux/internal/tmux"
)

func TestMergeTwoIdentical(t *testing.T) {
	r := []Rule{
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
	assert.Equal(t, w[0].Name, "editor")
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, w[0].Panes[0].Cmd, "echo test")
}

func TestMerge3_2identical(t *testing.T) {
	r := []Rule{
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
	assert.Equal(t, w[0].Name, "editor")
	assert.Equal(t, w[1].Name, "logs")
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, w[0].Panes[0].Cmd, "echo test")
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, w[1].Panes[0].Cmd, "echo logs")
}

func TestMerge3_2identical_different_order(t *testing.T) {
	r := []Rule{
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
	assert.Equal(t, w[0].Name, "logs")
	assert.Equal(t, w[1].Name, "editor")
	assert.Len(t, w[0].Panes, 1)
	assert.Equal(t, w[0].Panes[0].Cmd, "echo logs")
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, w[1].Panes[0].Cmd, "echo test")
}

func TestMerge3_MergePanes(t *testing.T) {
	r := []Rule{
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
	assert.Equal(t, w[0].Name, "editor")
	assert.Equal(t, w[1].Name, "logs")
	assert.Len(t, w[0].Panes, 2)
	assert.Equal(t, w[0].Panes[0].Cmd, "echo test")
	assert.Equal(t, w[0].Panes[1].Cmd, "echo another test")
	assert.Len(t, w[1].Panes, 1)
	assert.Equal(t, w[1].Panes[0].Cmd, "echo logs")
}
