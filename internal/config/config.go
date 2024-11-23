package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
	"path/filepath"
)

const APP_NAME = "automux"

type Config struct {
	Projects ProjectsConfig `toml:"projects"`
}

type ProjectsConfig struct {
	LookupDirs []string `toml:"lookup_dirs"`
}

func getConfigPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdgConfigHome, APP_NAME, "config.toml")
}

func (c *Config) Load() error {
	configPath := getConfigPath()

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read config file: %v", err)
	}

	if err := toml.Unmarshal(bytes, c); err != nil {
		return fmt.Errorf("could not unmarshal config: %v", err)
	}

	return nil
}
