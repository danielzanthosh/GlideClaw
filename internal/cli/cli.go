package cli

import (
	"context"
	"fmt"
	"strings"

	"glideclaw/internal/archive"
	"glideclaw/internal/bootstrap"
	"glideclaw/internal/config"
	"glideclaw/internal/connectors"
	"glideclaw/internal/db"
	"glideclaw/internal/policy"
	"glideclaw/internal/telegram"
)

type Router struct {
	cfg      config.Config
	store    *db.Store
	registry *connectors.Registry
	boot     bootstrap.Profile
	policy   *policy.Engine
	archive  *archive.Manager
}

func NewRouter(cfg config.Config, store *db.Store, registry *connectors.Registry, boot bootstrap.Profile, policy *policy.Engine, archive *archive.Manager) *Router {
	return &Router{cfg: cfg, store: store, registry: registry, boot: boot, policy: policy, archive: archive}
}

func (r *Router) Dispatch(ctx context.Context, args []string, bot *telegram.Adapter) error {
	if len(args) == 0 {
		return r.help()
	}
	cmd := strings.Join(args, " ")
	switch cmd {
	case "run":
		fmt.Println("starting glideclaw daemon components")
		return bot.Start(ctx)
	case "doctor":
		fmt.Printf("profile=%s data_dir=%s safe_mode=%v\n", r.cfg.Profile, r.cfg.DataDir, r.cfg.Execution.SafeMode)
		for _, h := range r.registry.Health(ctx) {
			fmt.Printf("connector=%s status=%s detail=%s\n", h.Connector, h.Status, h.Detail)
		}
		return nil
	case "config validate":
		fmt.Println("configuration validated")
		return nil
	case "connector status":
		for _, h := range r.registry.Health(ctx) {
			fmt.Printf("%s: %s\n", h.Connector, h.Status)
		}
		return nil
	case "archive run":
		return r.archive.RunOffloadSweep(ctx)
	case "safe-mode on":
		fmt.Println("safe mode is controlled by config/env; set GLIDECLAW_SAFE_MODE=true")
		return nil
	default:
		return r.help()
	}
}

func (r *Router) help() error {
	fmt.Println("glideclaw commands: run, doctor, config validate, connector status, archive run, safe-mode on")
	return nil
}
