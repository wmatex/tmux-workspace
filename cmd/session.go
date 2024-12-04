package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wmatex/automux/internal/projects"
	"github.com/wmatex/automux/internal/rules"
	"github.com/wmatex/automux/internal/tmux"
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

		projectDirs := viper.GetStringSlice("projects.lookup_dirs")

		allRules, err := rules.LoadFromConfig()
		if err != nil {
			log.Fatal(err)
		}

		p, err := projects.LoadAllProjects(projectDirs)
		if err != nil {
			log.Fatalf("cannot load all projects: %s\n", err)
		}

		sessions, err := tmux.GetActiveSessions()
		if err != nil {
			log.Fatalf("cannot get active sessions: %s\n", err)
		}

		p.MergeProjectsWithSessions(sessions)

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
