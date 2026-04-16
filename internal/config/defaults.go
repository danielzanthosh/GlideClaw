// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package config

func defaults() Config {
	return Config{
		Profile:   "lite",
		DataDir:   "/var/lib/glideclaw",
		Bootstrap: BootstrapConfig{Path: "BOOTSTRAP.md"},
		Database:  DatabaseConfig{Path: "state/glideclaw.db"},
		Telegram: TelegramConfig{
			PollSeconds: 2,
		},
		Archive: ArchiveConfig{
			HotDir:            "hot",
			RestoreCacheDir:   "restore-cache",
			OffloadAfterHours: 168,
			MaxHotMB:          512,
		},
		Execution: ExecutionConfig{
			SafeMode:                true,
			WorkspaceAllow:          []string{"/srv/projects", "/home"},
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
			SecretsDir:                "secrets",
		},
		Logging: LoggingConfig{Level: "info", Path: "logs/glideclaw.jsonl"},
	}
}
