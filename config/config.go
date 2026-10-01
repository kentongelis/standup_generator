package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the user's standup settings
type Config struct {
	GitHubUsername string   `mapstructure:"github_username"`
	GitHubToken    string   `mapstructure:"github_token"`
	GitEmail       string   `mapstructure:"git_email"` // falls back to git's global user.email
	Repos          []string `mapstructure:"repos"`
}

// Load reads ~/.standup.yaml, falling back to defaults where it can
func Load() (*Config, error) {
	v := viper.New()

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	v.AddConfigPath(home)
	v.SetConfigName("standup")
	v.SetConfigType("yaml")
	v.BindEnv("github_token", "GITHUB_TOKEN")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
		// no config file is fine, just use defaults/env
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// no email configured, ask git for the user's global one
	if cfg.GitEmail == "" {
		out, err := exec.Command("git", "config", "--global", "user.email").Output()
		if err == nil {
			cfg.GitEmail = strings.TrimSpace(string(out))
		}
	}
	return &cfg, nil
}
