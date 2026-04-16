package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Profile   string
	DataDir   string
	Bootstrap BootstrapConfig
	Database  DatabaseConfig
	Telegram  TelegramConfig
	Archive   ArchiveConfig
	Execution ExecutionConfig
	Logging   LoggingConfig
}

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
	SafeMode          bool
	WorkspaceAllow    []string
	Tier0Allow        []string
	Tier3DenyPrefix   []string
	DefaultTimeoutSec int
}

type LoggingConfig struct {
	Level string
	Path  string
}

func Load(path string) (Config, error) {
	cfg := defaults()

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
	if c.Bootstrap.Path == "" {
		c.Bootstrap.Path = filepath.Join(c.DataDir, "BOOTSTRAP.md")
	}
	return nil
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
