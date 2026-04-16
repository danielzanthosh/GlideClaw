// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package config

import "fmt"

func defaults(paths Paths) Config {
	return Config{
		Profile:    "lite",
		ConfigPath: paths.ConfigPath,
		DataDir:    paths.DataDir,
		Bootstrap:  BootstrapConfig{Path: paths.Bootstrap},
		Database:   DatabaseConfig{Path: "state/glideclaw.db"},
		Runtime:    RuntimeConfig{PIDFile: paths.PIDFile},
		Telegram: TelegramConfig{
			PollSeconds: 2,
		},
		Archive: ArchiveConfig{
			HotDir:            "cache/hot",
			RestoreCacheDir:   "cache/restore",
			OffloadAfterHours: 168,
			MaxHotMB:          512,
		},
		Execution: ExecutionConfig{
			SafeMode:                true,
			WorkspaceAllow:          []string{paths.DataDir},
			Tier0Allow:              []string{"pwd", "ls", "cat", "find", "du", "df", "git status", "git diff"},
			Tier3EscalationPrefixes: []string{"rm ", "git push", "git commit", "systemctl restart", "vercel --prod", "apt ", "dnf ", "yum "},
			HardBlockPrefixes:       []string{"sudo", "useradd", "usermod", "iptables", "nft", "mkfs", "fdisk", "parted", "rm -rf /"},
			DefaultTimeoutSec:       60,
		},
		Security: SecurityConfig{
			EscalationEnabled:         true,
			ElevationMode:             "time_window",
			ElevationWindowSeconds:    60,
			MaxAttempts:               5,
			LockoutSeconds:            60,
			RequireDoubleConfirmation: true,
			CriticalConfirmText:       "DELETE_PRODUCTION_DATA",
			AllowTier3InSafeMode:      false,
			SecretsDir:                paths.Secrets,
		},
		Logging: LoggingConfig{Level: "info", Path: "logs/glideclaw.jsonl"},
	}
}

func DefaultConfigYAML(cfg Config) string {
	return fmt.Sprintf(`profile: %s
data_dir: %s

execution:
  safe_mode: %t

security:
  escalation_enabled: %t
  elevation_mode: %s
  elevation_window_seconds: %d
  max_attempts: %d
  lockout_seconds: %d
  secrets_dir: %s

archive:
  hot_dir: %s
  restore_cache_dir: %s

connectors:
  google_drive:
    enabled: false
  google_gmail:
    enabled: false
  google_calendar:
    enabled: false
  github:
    enabled: false
  vercel:
    enabled: false
`, cfg.Profile, cfg.DataDir, cfg.Execution.SafeMode, cfg.Security.EscalationEnabled, cfg.Security.ElevationMode, cfg.Security.ElevationWindowSeconds, cfg.Security.MaxAttempts, cfg.Security.LockoutSeconds, cfg.Security.SecretsDir, cfg.Archive.HotDir, cfg.Archive.RestoreCacheDir)
}
