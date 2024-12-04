package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wmatex/automux/internal/fzf"
	"github.com/wmatex/automux/internal/projects"
	"github.com/wmatex/automux/internal/rules"
	"github.com/wmatex/automux/internal/tmux"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const APP_NAME = "automux"

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

		merged := p.MergeProjectsWithSessions(sessions)

		projectName, err := fzf.ProjectPick(merged)
		if err != nil {
			log.Fatalf("fzf failed: %s\n", err)
		}

		if projectName == "" {
			os.Exit(0)
		}

		project, ok := p.Map[projectName]
		if !ok {
			err, _ = tmux.NewSession(projectName, "")
			if err != nil {
				log.Fatalf("cannot create new tmux session '%s': %s\n", projectName, err)
			}
		} else if !project.Running {
			err, _ = tmux.NewSession(project.Name, project.Path)
			if err != nil {
				log.Fatalf("cannot create new tmux session '%s': %s\n", project.Name, err)
			}

			valid := allRules.GetSatisfied(project)
			windows := rules.MergeWindows(valid)

			err = rules.SetupHooks(project, valid)
			if err != nil {
				log.Fatalf("cannot setup hooks for project: %s\n", err)
			}

			err = tmux.CreateWindowsForProject(project.Name, project.Path, windows)
			if err != nil {
				log.Fatalf("cannot create windows: %s\n", err)
			}
		}

		err, _ = tmux.SwitchToSession(projectName)
		if err != nil {
			log.Fatalf("cannot switch to session '%s': %s\n", projectName, err)
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

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/automux/config.toml)")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
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

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName(fmt.Sprintf("%s.yaml", APP_NAME))
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
