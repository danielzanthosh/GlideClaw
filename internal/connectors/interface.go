package connectors

import "context"

type Health struct {
	Connector string
	Status    string
	Detail    string
}

type Connector interface {
	Name() string
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Health(ctx context.Context) Health
}

type Noop struct{ name string }

func NewNoop(name string) *Noop               { return &Noop{name: name} }
func (n *Noop) Name() string                  { return n.name }
func (n *Noop) Enable(context.Context) error  { return nil }
func (n *Noop) Disable(context.Context) error { return nil }
func (n *Noop) Health(context.Context) Health {
	return Health{Connector: n.name, Status: "not_configured", Detail: "skeleton connector"}
}
