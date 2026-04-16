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
			SafeMode:          true,
			WorkspaceAllow:    []string{"/srv/projects", "/home"},
			Tier0Allow:        []string{"pwd", "ls", "cat", "find", "du", "df", "git status", "git diff"},
			Tier3DenyPrefix:   []string{"sudo", "useradd", "usermod", "iptables", "nft", "mkfs", "fdisk", "parted", "rm -rf /"},
			DefaultTimeoutSec: 60,
		},
		Logging: LoggingConfig{Level: "info", Path: "logs/glideclaw.jsonl"},
	}
}
