package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(path string, mode string) (*Config, error) {
	if mode == "production" {
		path = ProductionConfigPath
	} else {
		path = DevelopmentConfigPath
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("TOML decode failed: %w", err)
	}
	applyDefaults(cfg)
	return cfg, nil
}

func applyDefaults(c *Config) {
	// Server
	if c.Server.Bind == "" {
		c.Server.Bind = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8443
	}
	// Logging
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.File == "" {
		c.Log.File = "/var/log/kiwipanel.log"
	}

	// Paths
	if c.Paths.DataDir == "" {
		c.Paths.DataDir = "/opt/kiwipanel/data"
	}
	if c.Paths.TempDir == "" {
		c.Paths.TempDir = "/opt/kiwipanel/tmp"
	}
	if c.Paths.BinaryPath == "" {
		c.Paths.BinaryPath = "/opt/kiwipanel/bin/kiwipanel"
	}
}
