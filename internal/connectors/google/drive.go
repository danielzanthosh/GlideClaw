// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package google

import (
	"context"

	"glideclaw/internal/connectors"
)

type Drive struct{}

func (d *Drive) Name() string                  { return "google_drive" }
func (d *Drive) Enable(context.Context) error  { return nil }
func (d *Drive) Disable(context.Context) error { return nil }
func (d *Drive) Health(context.Context) connectors.Health {
	return connectors.Health{Connector: d.Name(), Status: "todo", Detail: "device oauth + archive api"}
}
