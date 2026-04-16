// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package telegram

import (
	"context"

	"glideclaw/internal/archive"
	"glideclaw/internal/config"
	"glideclaw/internal/db"
	"glideclaw/internal/policy"
)

type Adapter struct {
	cfg    config.TelegramConfig
	store  *db.Store
	policy *policy.Engine
	arch   *archive.Manager
}

func NewAdapter(cfg config.TelegramConfig, store *db.Store, policy *policy.Engine, arch *archive.Manager) *Adapter {
	return &Adapter{cfg: cfg, store: store, policy: policy, arch: arch}
}

func (a *Adapter) Start(ctx context.Context) error {
	_ = ctx
	// Skeleton: implement long-polling, pairing checks, DM/group routing, and attachment intake.
	// Tier3 UX note: password challenge replies must be ephemeral and never persisted as message history.
	return nil
}
