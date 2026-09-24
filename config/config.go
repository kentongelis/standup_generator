package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	GitHubUsername string   `mapstructure:"github_username"`
	GitHubToken    string   `mapstructure:"github_token"`
	Repos          []string `mapstructure:"repos"`
}

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
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}
