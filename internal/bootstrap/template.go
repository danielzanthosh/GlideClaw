// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package bootstrap

import (
	"os"
	"path/filepath"
)

const defaultTemplate = `# GlideClaw Bootstrap Profile

## Identity
- Personal terminal-first operator.

## Preferences
- Keep responses concise and actionable.

## Allowed autonomous actions
- summarize logs

## Confirmation-required actions
- package install
- git push

## Blocked actions
- sudo

## Security mode
- strict
`

func EnsureProfile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if data, err := os.ReadFile("BOOTSTRAP.example.md"); err == nil {
		return os.WriteFile(path, data, 0o644)
	}
	return os.WriteFile(path, []byte(defaultTemplate), 0o644)
}
