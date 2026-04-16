// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"glideclaw/internal/bootstrap"
	"glideclaw/internal/config"
	"glideclaw/internal/db"
)

type Result struct {
	ConfigPath         string
	DataDir            string
	DBPath             string
	BootstrapPath      string
	SecretsPlaceholder string
	PIDPath            string
}

func Init(cfg config.Config) (Result, error) {
	dirs := []string{
		filepath.Dir(cfg.ConfigPath),
		cfg.DataDir,
		filepath.Dir(cfg.Database.Path),
		cfg.Archive.HotDir,
		cfg.Archive.RestoreCacheDir,
		filepath.Dir(cfg.Logging.Path),
		cfg.Security.SecretsDir,
		filepath.Dir(cfg.Runtime.PIDFile),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return Result{}, err
		}
	}

	if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
		if err := config.WriteDefaultConfig(cfg.ConfigPath, cfg); err != nil {
			return Result{}, err
		}
	}

	store, err := db.OpenAndMigrate(cfg.Database.Path)
	if err != nil {
		return Result{}, err
	}
	_ = store.DB.Close()

	secretsPlaceholder := filepath.Join(cfg.Security.SecretsDir, "README.txt")
	if _, err := os.Stat(secretsPlaceholder); os.IsNotExist(err) {
		_ = os.WriteFile(secretsPlaceholder, []byte("GlideClaw secret placeholder. Set escalation password via: glideclaw security set-password\n"), 0o600)
	}

	if err := bootstrap.EnsureProfile(cfg.Bootstrap.Path); err != nil {
		return Result{}, err
	}

	return Result{
		ConfigPath:         cfg.ConfigPath,
		DataDir:            cfg.DataDir,
		DBPath:             cfg.Database.Path,
		BootstrapPath:      cfg.Bootstrap.Path,
		SecretsPlaceholder: secretsPlaceholder,
		PIDPath:            cfg.Runtime.PIDFile,
	}, nil
}

func NextSteps(r Result) string {
	return fmt.Sprintf(`Initialization complete.

Next steps:
1) Review config: %s
2) Customize bootstrap profile: %s
3) Set escalation password: glideclaw security set-password
4) Start daemon: glideclaw run
`, r.ConfigPath, r.BootstrapPath)
}
