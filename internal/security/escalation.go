// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package security

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"glideclaw/internal/config"
)

type lockoutState struct {
	FailedAttempts int       `json:"failed_attempts"`
	LockedUntil    time.Time `json:"locked_until"`
	ElevatedUntil  time.Time `json:"elevated_until"`
}

type EscalationManager struct {
	mu         sync.Mutex
	cfg        config.SecurityConfig
	state      lockoutState
	path       string
	secretsDir string
}

func NewEscalationManager(cfg config.SecurityConfig) (*EscalationManager, error) {
	if err := os.MkdirAll(cfg.SecretsDir, 0o700); err != nil {
		return nil, err
	}
	m := &EscalationManager{cfg: cfg, path: filepath.Join(cfg.SecretsDir, "escalation_state.json"), secretsDir: cfg.SecretsDir}
	_ = m.load()
	return m, nil
}

func (m *EscalationManager) load() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.state)
}

func (m *EscalationManager) save() error {
	data, err := json.Marshal(m.state)
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o600)
}

func (m *EscalationManager) IsLocked(now time.Time) (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.LockedUntil.After(now) {
		return true, m.state.LockedUntil.Sub(now)
	}
	return false, 0
}

func (m *EscalationManager) RecordFailure(now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.FailedAttempts++
	if m.state.FailedAttempts >= m.cfg.MaxAttempts {
		m.state.LockedUntil = now.Add(time.Duration(m.cfg.LockoutSeconds) * time.Second)
		m.state.FailedAttempts = 0
	}
	return m.save()
}

func (m *EscalationManager) RecordSuccess(now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.FailedAttempts = 0
	m.state.LockedUntil = time.Time{}
	if m.cfg.ElevationMode == "time_window" {
		m.state.ElevatedUntil = now.Add(time.Duration(m.cfg.ElevationWindowSeconds) * time.Second)
	} else {
		m.state.ElevatedUntil = now.Add(2 * time.Second)
	}
	return m.save()
}

func (m *EscalationManager) IsElevated(now time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state.ElevatedUntil.After(now)
}

func (m *EscalationManager) ConsumeSingleUse(now time.Time) {
	if m.cfg.ElevationMode != "single" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.ElevatedUntil.After(now) {
		m.state.ElevatedUntil = time.Time{}
		_ = m.save()
	}
}

func (m *EscalationManager) ResetLockout() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.FailedAttempts = 0
	m.state.LockedUntil = time.Time{}
	return m.save()
}

func (m *EscalationManager) Status(now time.Time) string {
	locked, remaining := m.IsLocked(now)
	if locked {
		return "locked for " + remaining.Round(time.Second).String()
	}
	if m.IsElevated(now) {
		return "elevated"
	}
	if !PasswordConfigured(m.cfg.SecretsDir) {
		return "password_not_set"
	}
	return "ready"
}

func (m *EscalationManager) RequireEnabled() error {
	if !m.cfg.EscalationEnabled {
		return errors.New("escalation disabled by config")
	}
	return nil
}

func (m *EscalationManager) SecretsDir() string {
	return m.secretsDir
}
