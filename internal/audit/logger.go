// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"glideclaw/internal/db"
)

type Logger struct {
	store *db.Store
}

func NewLogger(store *db.Store) *Logger {
	return &Logger{store: store}
}

type Tier3Event struct {
	Command     string
	Source      string
	RequestedBy string
	Result      string
	ExecResult  string
}

func (l *Logger) LogTier3Attempt(ctx context.Context, ev Tier3Event) error {
	details, _ := json.Marshal(map[string]string{
		"source":      ev.Source,
		"exec_result": ev.ExecResult,
	})
	_, err := l.store.DB.ExecContext(ctx, `
		INSERT INTO audit_log (id, actor, action, target, connector, outcome, details_json, created_at)
		VALUES (?, ?, 'tier3_attempt', ?, '', ?, ?, ?)
	`, uuid.NewString(), ev.RequestedBy, ev.Command, ev.Result, string(details), time.Now().UTC())
	return err
}
