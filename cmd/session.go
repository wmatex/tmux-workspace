package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wmatex/tmux-workspace/internal/rules"
	"github.com/wmatex/tmux-workspace/internal/tmux"
)

var projectName string
var lifecycle string

var sessionCmd = &cobra.Command{
	Use:       "session [start|end|kill]",
	ValidArgs: []string{"start", "end", "kill"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			os.Exit(0)
		}

		p, allRules := initProjectsAndRules()

		project, ok := p.Map[projectName]
		if ok {
			if args[0] == "kill" {
				err, _ := tmux.KillSession(project.Name)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				var lifecycle uint8
				switch args[0] {
				case "start":
					lifecycle = rules.START
				case "end":
					lifecycle = rules.END
				}

				valid := allRules.GetSatisfied(project)
				err := rules.RunHooks(project, lifecycle, valid)
				if err != nil {
					os.Exit(1)
				}
			}

		}
	},
}

func init() {
	sessionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name (required)")
	sessionCmd.MarkFlagRequired("project")
}
