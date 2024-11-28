package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wmatex/automux/internal/cmd_exec/fzf"
	"github.com/wmatex/automux/internal/cmd_exec/tmux"
	"github.com/wmatex/automux/internal/projects"

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

		project := p.Map[projectName]

		if !project.Running {
			err = tmux.NewSession(project.Name, project.Path)
			if err != nil {
				log.Fatalf("cannot create new tmux session '%s': %s\n", project.Name, err)
			}
		}

		err = tmux.SwitchToSession(project.Name)
		if err != nil {
			log.Fatalf("cannot switch to session '%s': %s\n", project.Name, err)
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
		viper.SetConfigType("toml")
		viper.SetConfigName("config.toml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
