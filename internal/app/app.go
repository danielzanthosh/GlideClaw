// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package app

import (
	"context"

	"glideclaw/internal/archive"
	"glideclaw/internal/audit"
	"glideclaw/internal/bootstrap"
	"glideclaw/internal/cli"
	"glideclaw/internal/config"
	"glideclaw/internal/connectors"
	"glideclaw/internal/db"
	"glideclaw/internal/executor"
	"glideclaw/internal/policy"
	"glideclaw/internal/security"
	"glideclaw/internal/telegram"
)

type App struct {
	cfg       config.Config
	db        *db.Store
	bootstrap bootstrap.Profile
	policy    *policy.Engine
	archive   *archive.Manager
	bot       *telegram.Adapter
	cli       *cli.Router
}

func New() (*App, error) {
	cfg, err := config.Load("")
	if err != nil {
		return nil, err
	}

	store, err := db.OpenAndMigrate(cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	bp, err := bootstrap.Load(cfg.Bootstrap.Path)
	if err != nil {
		return nil, err
	}

	registry := connectors.NewRegistry()
	registry.Register(connectors.NewNoop("google_drive"))
	registry.Register(connectors.NewNoop("google_gmail"))
	registry.Register(connectors.NewNoop("google_calendar"))
	registry.Register(connectors.NewNoop("github"))
	registry.Register(connectors.NewNoop("vercel"))

	engine := policy.NewEngine(cfg.Execution, bp)
	archiver := archive.NewManager(cfg.Archive, store, registry)
	esc, err := security.NewEscalationManager(cfg.Security)
	if err != nil {
		return nil, err
	}
	auditLog := audit.NewLogger(store)
	runner := executor.NewRunner(engine, esc, auditLog, executor.SecurityConfig{
		RequireDoubleConfirmation: cfg.Security.RequireDoubleConfirmation,
		CriticalConfirmText:       cfg.Security.CriticalConfirmText,
		AllowTier3InSafeMode:      cfg.Security.AllowTier3InSafeMode,
	})
	bot := telegram.NewAdapter(cfg.Telegram, store, engine, archiver)
	router := cli.NewRouter(cfg, store, registry, bp, engine, archiver, esc, runner)

	return &App{cfg: cfg, db: store, bootstrap: bp, policy: engine, archive: archiver, bot: bot, cli: router}, nil
}

func (a *App) Run(ctx context.Context, args []string) error {
	return a.cli.Dispatch(ctx, args, a.bot)
}
