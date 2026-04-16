// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"glideclaw/internal/archive"
	"glideclaw/internal/bootstrap"
	"glideclaw/internal/config"
	"glideclaw/internal/connectors"
	"glideclaw/internal/db"
	"glideclaw/internal/executor"
	"glideclaw/internal/policy"
	"glideclaw/internal/security"
	"glideclaw/internal/setup"
	"glideclaw/internal/telegram"
)

type Router struct {
	cfg        config.Config
	store      *db.Store
	registry   *connectors.Registry
	boot       bootstrap.Profile
	policy     *policy.Engine
	archive    *archive.Manager
	escalation *security.EscalationManager
	runner     *executor.Runner
}

func NewRouter(cfg config.Config, store *db.Store, registry *connectors.Registry, boot bootstrap.Profile, policy *policy.Engine, archive *archive.Manager, escalation *security.EscalationManager, runner *executor.Runner) *Router {
	return &Router{cfg: cfg, store: store, registry: registry, boot: boot, policy: policy, archive: archive, escalation: escalation, runner: runner}
}

func (r *Router) Dispatch(ctx context.Context, args []string, bot *telegram.Adapter) error {
	if len(args) == 0 {
		return r.help()
	}
	cmd := strings.Join(args, " ")
	switch {
	case cmd == "init":
		res, err := setup.Init(r.cfg)
		if err != nil {
			return err
		}
		fmt.Println(setup.NextSteps(res))
		return nil
	case cmd == "status":
		return r.status(ctx)
	case cmd == "run":
		return r.runDaemon(ctx, bot)
	case cmd == "doctor":
		fmt.Printf("profile=%s data_dir=%s safe_mode=%v escalation=%s\n", r.cfg.Profile, r.cfg.DataDir, r.cfg.Execution.SafeMode, r.escalation.Status(time.Now()))
		for _, h := range r.registry.Health(ctx) {
			fmt.Printf("connector=%s status=%s detail=%s\n", h.Connector, h.Status, h.Detail)
		}
		return nil
	case cmd == "config validate":
		fmt.Println("configuration validated")
		return nil
	case cmd == "connector status":
		for _, h := range r.registry.Health(ctx) {
			fmt.Printf("%s: %s\n", h.Connector, h.Status)
		}
		return nil
	case cmd == "archive run":
		return r.archive.RunOffloadSweep(ctx)
	case cmd == "security set-password":
		return setPassword(r.cfg.Security.SecretsDir)
	case cmd == "security change-password":
		return changePassword(r.cfg.Security.SecretsDir)
	case cmd == "security status":
		fmt.Println(r.escalation.Status(time.Now()))
		return nil
	case cmd == "security reset-lockout":
		if err := r.escalation.ResetLockout(); err != nil {
			return err
		}
		fmt.Println("lockout reset")
		return nil
	case strings.HasPrefix(cmd, "exec "):
		raw := strings.TrimSpace(strings.TrimPrefix(cmd, "exec "))
		overrideSafe := false
		if strings.HasPrefix(raw, "--override-safe ") {
			overrideSafe = true
			raw = strings.TrimSpace(strings.TrimPrefix(raw, "--override-safe "))
		}
		return r.runner.Execute(ctx, raw, "terminal", "local_user", overrideSafe)
	default:
		return r.help()
	}
}

func (r *Router) runDaemon(ctx context.Context, bot *telegram.Adapter) error {
	if _, err := setup.Init(r.cfg); err != nil {
		return fmt.Errorf("runtime initialization failed: %w", err)
	}
	if err := writePID(r.cfg.Runtime.PIDFile); err != nil {
		return err
	}
	defer os.Remove(r.cfg.Runtime.PIDFile)

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("[glideclaw] starting daemon (profile=%s, data_dir=%s)\n", r.cfg.Profile, r.cfg.DataDir)
	if err := bot.Start(sigCtx); err != nil {
		fmt.Printf("[glideclaw] telegram adapter exited: %v\n", err)
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigCtx.Done():
			fmt.Println("[glideclaw] shutdown signal received")
			fmt.Println("[glideclaw] daemon stopped")
			return nil
		case <-ticker.C:
			fmt.Println("[glideclaw] heartbeat")
		}
	}
}

func (r *Router) status(ctx context.Context) error {
	bootstrapExists := fileExists(r.cfg.Bootstrap.Path)
	secretsExists := fileExists(r.cfg.Security.SecretsDir + "/escalation_password.json")
	pidExists := fileExists(r.cfg.Runtime.PIDFile)
	fmt.Printf("profile: %s\n", r.cfg.Profile)
	fmt.Printf("safe_mode: %v\n", r.cfg.Execution.SafeMode)
	fmt.Printf("data_dir: %s\n", r.cfg.DataDir)
	fmt.Printf("config_path: %s\n", r.cfg.ConfigPath)
	fmt.Printf("db_path: %s\n", r.cfg.Database.Path)
	fmt.Printf("bootstrap_exists: %v\n", bootstrapExists)
	fmt.Printf("secrets_exists: %v\n", secretsExists)
	fmt.Printf("daemon_pid_exists: %v\n", pidExists)
	fmt.Println("connectors:")
	for _, h := range r.registry.Health(ctx) {
		fmt.Printf("  - %s: %s\n", h.Connector, h.Status)
	}
	return nil
}

func (r *Router) help() error {
	fmt.Println("glideclaw commands: init, status, run, doctor, config validate, connector status, archive run, security [set-password|change-password|status|reset-lockout], exec [--override-safe] <cmd>")
	return nil
}

func setPassword(secretsDir string) error {
	fmt.Print("New escalation password: ")
	first, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}
	fmt.Print("Confirm escalation password: ")
	second, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}
	if string(first) != string(second) {
		return fmt.Errorf("passwords do not match")
	}
	if err := security.SetPassword(secretsDir, string(first)); err != nil {
		return err
	}
	fmt.Println("escalation password set")
	return nil
}

func changePassword(secretsDir string) error {
	if !security.PasswordConfigured(secretsDir) {
		return fmt.Errorf("password not configured")
	}
	fmt.Print("Current escalation password: ")
	current, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}
	ok, err := security.VerifyPassword(secretsDir, string(current))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invalid current password")
	}
	return setPassword(secretsDir)
}

func writePID(path string) error {
	if err := os.MkdirAll(filepathDir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644)
}

func filepathDir(path string) string {
	i := strings.LastIndex(path, "/")
	if i <= 0 {
		return "."
	}
	return path[:i]
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
