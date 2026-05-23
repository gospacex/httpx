package httpx

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Adapter string `mapstructure:"adapter"`
	Port    int    `mapstructure:"port"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.Adapter == "" {
		return nil, fmt.Errorf("adapter is required in config")
	}
	if cfg.Port <= 0 {
		cfg.Port = 8080
	}

	return &cfg, nil
}