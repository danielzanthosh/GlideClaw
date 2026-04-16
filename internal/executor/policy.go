// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package executor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"glideclaw/internal/audit"
	"glideclaw/internal/policy"
	"glideclaw/internal/security"
)

type SecurityConfig struct {
	RequireDoubleConfirmation bool
	CriticalConfirmText       string
	AllowTier3InSafeMode      bool
}

type Runner struct {
	policy     *policy.Engine
	escalation *security.EscalationManager
	audit      *audit.Logger
	security   SecurityConfig
}

func NewRunner(pol *policy.Engine, esc *security.EscalationManager, aud *audit.Logger, sec SecurityConfig) *Runner {
	return &Runner{policy: pol, escalation: esc, audit: aud, security: sec}
}

func (r *Runner) Execute(ctx context.Context, command string, source string, actor string, overrideSafe bool) error {
	decision := r.policy.Evaluate(command, "")
	if decision.HardBlocked {
		_ = r.audit.LogTier3Attempt(ctx, audit.Tier3Event{Command: command, Source: source, RequestedBy: actor, Result: "denied", ExecResult: decision.Reason})
		return errors.New(decision.Reason)
	}

	if decision.Tier == policy.Tier3 {
		if r.policy.SafeModeOn() && !overrideSafe {
			_ = r.audit.LogTier3Attempt(ctx, audit.Tier3Event{Command: command, Source: source, RequestedBy: actor, Result: "denied", ExecResult: "safe mode blocks tier3"})
			return errors.New("safe mode blocks tier3 unless --override-safe is used")
		}
		if r.policy.SafeModeOn() && overrideSafe && !r.security.AllowTier3InSafeMode {
			return errors.New("tier3 safe-mode override disabled by policy")
		}
		if err := r.requireEscalation(command); err != nil {
			_ = r.audit.LogTier3Attempt(ctx, audit.Tier3Event{Command: command, Source: source, RequestedBy: actor, Result: "denied", ExecResult: err.Error()})
			return err
		}
	}

	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Run()
	if decision.Tier == policy.Tier3 {
		res := "ok"
		if err != nil {
			res = err.Error()
		}
		_ = r.audit.LogTier3Attempt(ctx, audit.Tier3Event{Command: command, Source: source, RequestedBy: actor, Result: "approved", ExecResult: res})
		r.escalation.ConsumeSingleUse(time.Now())
	}
	return err
}

func (r *Runner) requireEscalation(command string) error {
	if err := r.escalation.RequireEnabled(); err != nil {
		return err
	}
	locked, remaining := r.escalation.IsLocked(time.Now())
	if locked {
		return fmt.Errorf("locked out, try again in %s", remaining.Round(time.Second))
	}
	if r.escalation.IsElevated(time.Now()) {
		return nil
	}

	fmt.Print("Escalation password: ")
	pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}
	ok, err := security.VerifyPassword(r.escalation.SecretsDir(), string(pwBytes))
	if err != nil {
		return err
	}
	if !ok {
		_ = r.escalation.RecordFailure(time.Now())
		return errors.New("invalid password")
	}

	if r.security.RequireDoubleConfirmation && isCritical(command) {
		fmt.Printf("Type: %s\n> ", r.security.CriticalConfirmText)
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) != r.security.CriticalConfirmText {
			_ = r.escalation.RecordFailure(time.Now())
			return errors.New("double confirmation failed")
		}
	}
	return r.escalation.RecordSuccess(time.Now())
}

func isCritical(command string) bool {
	critical := []string{"rm ", "DROP ", "truncate ", "systemctl stop", "vercel --prod"}
	for _, p := range critical {
		if strings.HasPrefix(command, p) || strings.Contains(command, p) {
			return true
		}
	}
	return false
}
