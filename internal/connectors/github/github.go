// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package github

import (
	"context"

	"glideclaw/internal/connectors"
)

type Connector struct{}

func (c *Connector) Name() string                  { return "github" }
func (c *Connector) Enable(context.Context) error  { return nil }
func (c *Connector) Disable(context.Context) error { return nil }
func (c *Connector) Health(context.Context) connectors.Health {
	return connectors.Health{Connector: c.Name(), Status: "todo", Detail: "fine-grained PAT/app"}
}
