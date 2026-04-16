// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package policy

import (
	"errors"
	"path/filepath"
	"strings"

	"glideclaw/internal/bootstrap"
	"glideclaw/internal/config"
)

type Tier int

const (
	Tier0 Tier = iota
	Tier1
	Tier2
	Tier3
)

type Decision struct {
	Tier        Tier
	Allowed     bool
	HardBlocked bool
	Reason      string
}

type Engine struct {
	exec config.ExecutionConfig
	boot bootstrap.Profile
}

func NewEngine(exec config.ExecutionConfig, boot bootstrap.Profile) *Engine {
	return &Engine{exec: exec, boot: boot}
}

func (e *Engine) Evaluate(command string, workdir string) Decision {
	command = strings.TrimSpace(command)
	if command == "" {
		return Decision{Tier: Tier3, Allowed: false, HardBlocked: true, Reason: "empty command"}
	}
	for _, b := range e.exec.HardBlockPrefixes {
		if strings.HasPrefix(command, b) {
			return Decision{Tier: Tier3, Allowed: false, HardBlocked: true, Reason: "blocked by hard denylist"}
		}
	}
	for _, b := range e.boot.BlockedActions {
		if strings.Contains(command, b) {
			return Decision{Tier: Tier3, Allowed: false, HardBlocked: true, Reason: "blocked by BOOTSTRAP"}
		}
	}
	for _, a := range e.exec.Tier0Allow {
		if command == a || strings.HasPrefix(command, a+" ") {
			return Decision{Tier: Tier0, Allowed: true, Reason: "tier0 allow"}
		}
	}
	for _, t3 := range e.exec.Tier3EscalationPrefixes {
		if strings.HasPrefix(command, t3) {
			return Decision{Tier: Tier3, Allowed: true, Reason: "tier3 escalation required"}
		}
	}
	if e.exec.SafeMode {
		return Decision{Tier: Tier2, Allowed: false, Reason: "safe mode requires explicit approval for non-tier0"}
	}
	return Decision{Tier: Tier1, Allowed: true, Reason: "policy allowed"}
}

func (e *Engine) ValidateWorkdir(path string) error {
	if path == "" {
		return nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for _, root := range e.exec.WorkspaceAllow {
		if strings.HasPrefix(abs, root) {
			return nil
		}
	}
	return errors.New("workdir is outside approved workspaces")
}

func (e *Engine) SafeModeOn() bool {
	return e.exec.SafeMode
}
