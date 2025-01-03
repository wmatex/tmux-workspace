package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wmatex/tmux-workspace/internal/fzf"
	"github.com/wmatex/tmux-workspace/internal/projects"
	"github.com/wmatex/tmux-workspace/internal/rules"
	"github.com/wmatex/tmux-workspace/internal/tmux"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const APP_NAME = "tmux-workspace"
const DEFAULT_LAYOUT = "main-vertical"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   APP_NAME,
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		windowLayout := viper.GetString("projects.layout")
		p, allRules := initProjectsAndRules()

		merged := p.GetNotActiveProjects()
		projectPick, err := fzf.ProjectPick(merged)
		if err != nil {
			log.Fatalf("fzf failed: %s\n", err)
		}

		if projectPick.ProjectName == "" {
			os.Exit(0)
		}

		project, ok := p.Map[projectPick.ProjectName]
		if !ok {
			err, _ = tmux.NewSession(projectPick.ProjectName, "")
			if err != nil {
				log.Fatalf("cannot create new tmux session '%s': %s\n", projectPick.ProjectName, err)
			}
		} else if !project.Running {
			valid := allRules.GetSatisfied(project)
			windows := rules.MergeWindows(valid)

			skipStartHook := false
			if projectPick.Action == fzf.SKIP_START_HOOK {
				skipStartHook = true
			}
			err = rules.SetupHooks(project, valid, skipStartHook)
			if err != nil {
				log.Fatalf("cannot setup hooks for project: %s\n", err)
			}

			err, _ = tmux.NewSession(project.Name, project.Path)
			if err != nil {
				log.Fatalf("cannot create new tmux session '%s': %s\n", project.Name, err)
			}

			err = tmux.CreateWindowsForProject(project.Name, project.Path, windowLayout, windows)
			if err != nil {
				log.Fatalf("cannot create windows: %s\n", err)
			}
		}

		err, _ = tmux.SwitchToSession(projectPick.ProjectName)
		if err != nil {
			log.Fatalf("cannot switch to session '%s': %s\n", projectPick.ProjectName, err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initProjectsAndRules() (*projects.Projects, *rules.Rules) {
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

	return p, allRules
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/tmux-workspace/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			xdgConfigHome = filepath.Join(home, ".config")
		}
		configPath := filepath.Join(xdgConfigHome, APP_NAME)

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.ReadInConfig()
}
