package google

import (
	"context"

	"glideclaw/internal/connectors"
)

type Gmail struct{}

func (g *Gmail) Name() string                  { return "google_gmail" }
func (g *Gmail) Enable(context.Context) error  { return nil }
func (g *Gmail) Disable(context.Context) error { return nil }
func (g *Gmail) Health(context.Context) connectors.Health {
	return connectors.Health{Connector: g.Name(), Status: "todo", Detail: "scoped gmail read/send"}
}
