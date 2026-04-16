package google

import (
	"context"

	"glideclaw/internal/connectors"
)

type Calendar struct{}

func (c *Calendar) Name() string                  { return "google_calendar" }
func (c *Calendar) Enable(context.Context) error  { return nil }
func (c *Calendar) Disable(context.Context) error { return nil }
func (c *Calendar) Health(context.Context) connectors.Health {
	return connectors.Health{Connector: c.Name(), Status: "todo", Detail: "event read/write"}
}
