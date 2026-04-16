// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profile    string
	ConfigPath string
	DataDir    string
	Bootstrap  BootstrapConfig
	Database   DatabaseConfig
	Telegram   TelegramConfig
	Archive    ArchiveConfig
	Execution  ExecutionConfig
	Security   SecurityConfig
	Logging    LoggingConfig
	Runtime    RuntimeConfig
}

type RuntimeConfig struct{ PIDFile string }

type BootstrapConfig struct{ Path string }
type DatabaseConfig struct{ Path string }

type TelegramConfig struct {
	BotToken    string
	PollSeconds int
	AllowedChat []int64
}

type ArchiveConfig struct {
	HotDir            string
	RestoreCacheDir   string
	DryRun            bool
	OffloadAfterHours int
	MaxHotMB          int
}

type ExecutionConfig struct {
	SafeMode                bool
	WorkspaceAllow          []string
	Tier0Allow              []string
	Tier3EscalationPrefixes []string
	HardBlockPrefixes       []string
	DefaultTimeoutSec       int
}

type SecurityConfig struct {
	EscalationEnabled         bool
	ElevationMode             string
	ElevationWindowSeconds    int
	MaxAttempts               int
	LockoutSeconds            int
	RequireDoubleConfirmation bool
	CriticalConfirmText       string
	AllowTier3InSafeMode      bool
	SecretsDir                string
}

type LoggingConfig struct {
	Level string
	Path  string
}

type diskConfig struct {
	Profile   string `yaml:"profile"`
	DataDir   string `yaml:"data_dir"`
	Execution struct {
		SafeMode bool `yaml:"safe_mode"`
	} `yaml:"execution"`
	Security struct {
		EscalationEnabled      *bool  `yaml:"escalation_enabled"`
		ElevationMode          string `yaml:"elevation_mode"`
		ElevationWindowSeconds int    `yaml:"elevation_window_seconds"`
		MaxAttempts            int    `yaml:"max_attempts"`
		LockoutSeconds         int    `yaml:"lockout_seconds"`
		SecretsDir             string `yaml:"secrets_dir"`
	} `yaml:"security"`
}

func Load(path string) (Config, error) {
	paths := ResolvePaths()
	if path == "" {
		path = paths.ConfigPath
	}
	cfg := defaults(paths)
	cfg.ConfigPath = path

	if data, err := os.ReadFile(path); err == nil {
		var dc diskConfig
		if err := yaml.Unmarshal(data, &dc); err == nil {
			if dc.Profile != "" {
				cfg.Profile = dc.Profile
			}
			if dc.DataDir != "" {
				cfg.DataDir = dc.DataDir
			}
			cfg.Execution.SafeMode = dc.Execution.SafeMode
			if dc.Security.EscalationEnabled != nil {
				cfg.Security.EscalationEnabled = *dc.Security.EscalationEnabled
			}
			if dc.Security.ElevationMode != "" {
				cfg.Security.ElevationMode = dc.Security.ElevationMode
			}
			if dc.Security.ElevationWindowSeconds > 0 {
				cfg.Security.ElevationWindowSeconds = dc.Security.ElevationWindowSeconds
			}
			if dc.Security.MaxAttempts > 0 {
				cfg.Security.MaxAttempts = dc.Security.MaxAttempts
			}
			if dc.Security.LockoutSeconds > 0 {
				cfg.Security.LockoutSeconds = dc.Security.LockoutSeconds
			}
			if dc.Security.SecretsDir != "" {
				cfg.Security.SecretsDir = dc.Security.SecretsDir
			}
		}
	}

	if v := os.Getenv("GLIDECLAW_PROFILE"); v != "" {
		cfg.Profile = v
	}
	if v := os.Getenv("GLIDECLAW_DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("GLIDECLAW_DB_PATH"); v != "" {
		cfg.Database.Path = v
	}
	if v := os.Getenv("GLIDECLAW_BOOTSTRAP_PATH"); v != "" {
		cfg.Bootstrap.Path = v
	}
	if v := os.Getenv("GLIDECLAW_TELEGRAM_BOT_TOKEN"); v != "" {
		cfg.Telegram.BotToken = v
	}
	if v := os.Getenv("GLIDECLAW_SAFE_MODE"); v != "" {
		cfg.Execution.SafeMode = strings.EqualFold(v, "true")
	}
	if v := os.Getenv("GLIDECLAW_EXEC_TIMEOUT_SEC"); v != "" {
		n, _ := strconv.Atoi(v)
		if n > 0 {
			cfg.Execution.DefaultTimeoutSec = n
		}
	}

	if err := cfg.Normalize(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) Normalize() error {
	if c.DataDir == "" {
		return errors.New("data_dir cannot be empty")
	}
	c.Database.Path = expand(c.Database.Path, c.DataDir)
	c.Logging.Path = expand(c.Logging.Path, c.DataDir)
	c.Archive.HotDir = expand(c.Archive.HotDir, c.DataDir)
	c.Archive.RestoreCacheDir = expand(c.Archive.RestoreCacheDir, c.DataDir)
	c.Security.SecretsDir = expand(c.Security.SecretsDir, c.DataDir)
	c.Runtime.PIDFile = expand(c.Runtime.PIDFile, c.DataDir)
	if c.Bootstrap.Path == "" {
		c.Bootstrap.Path = filepath.Join(filepath.Dir(c.ConfigPath), "BOOTSTRAP.md")
	}
	if c.Security.ElevationMode != "single" && c.Security.ElevationMode != "time_window" {
		return errors.New("security.elevation_mode must be single or time_window")
	}
	return nil
}

func WriteDefaultConfig(path string, cfg Config) error {
	content := []byte(DefaultConfigYAML(cfg))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func expand(path, dataDir string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dataDir, path)
}
