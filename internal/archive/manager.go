// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package archive

import (
	"context"
	"fmt"

	"glideclaw/internal/config"
	"glideclaw/internal/connectors"
	"glideclaw/internal/db"
)

type Manager struct {
	cfg      config.ArchiveConfig
	store    *db.Store
	registry *connectors.Registry
}

func NewManager(cfg config.ArchiveConfig, store *db.Store, registry *connectors.Registry) *Manager {
	return &Manager{cfg: cfg, store: store, registry: registry}
}

func (m *Manager) RunOffloadSweep(ctx context.Context) error {
	_ = ctx
	if m.cfg.DryRun {
		return nil
	}
	// Placeholder for selecting cold candidates and uploading to Drive connector.
	return nil
}

func (m *Manager) Restore(ctx context.Context, objectID string) (string, error) {
	_ = ctx
	if objectID == "" {
		return "", fmt.Errorf("object id is required")
	}
	// Placeholder: query archive_objects, fetch from Drive, verify checksum, return local cached path.
	return "", nil
}
