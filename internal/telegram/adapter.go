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
	return nil
}
