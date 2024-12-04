package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wmatex/automux/internal/rules"
)

var projectName string
var lifecycle string

var sessionCmd = &cobra.Command{
	Use:       "session [start|end]",
	ValidArgs: []string{"start", "end"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			os.Exit(0)
		}

		p, allRules := initProjectsAndRules()

		project, ok := p.Map[projectName]
		var lifecycle uint8
		switch args[0] {
		case "start":
			lifecycle = rules.START
		case "end":
			lifecycle = rules.END
		}
		if ok {
			valid := allRules.GetSatisfied(project)
			rules.RunHooks(project, lifecycle, valid)
		}
	},
}

func init() {
	sessionCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name (required)")
	sessionCmd.MarkFlagRequired("project")
}
