// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	ConfigDir  string
	ConfigPath string
	DataDir    string
	Bootstrap  string
	Secrets    string
	PIDFile    string
}

func ResolvePaths() Paths {
	home, _ := os.UserHomeDir()
	isRoot := os.Geteuid() == 0
	if isRoot {
		return Paths{
			ConfigDir:  "/etc/glideclaw",
			ConfigPath: "/etc/glideclaw/config.yaml",
			DataDir:    "/var/lib/glideclaw",
			Bootstrap:  "/etc/glideclaw/BOOTSTRAP.md",
			Secrets:    "/var/lib/glideclaw/secrets",
			PIDFile:    "/var/lib/glideclaw/run/glideclaw.pid",
		}
	}
	cfgDir := filepath.Join(home, ".config", "glideclaw")
	dataDir := filepath.Join(home, ".local", "share", "glideclaw")
	return Paths{
		ConfigDir:  cfgDir,
		ConfigPath: filepath.Join(cfgDir, "config.yaml"),
		DataDir:    dataDir,
		Bootstrap:  filepath.Join(cfgDir, "BOOTSTRAP.md"),
		Secrets:    filepath.Join(dataDir, "secrets"),
		PIDFile:    filepath.Join(dataDir, "run", "glideclaw.pid"),
	}
}
