package connectors

import (
	"context"
	"sort"
	"sync"
)

type Registry struct {
	mu   sync.RWMutex
	impl map[string]Connector
}

func NewRegistry() *Registry {
	return &Registry{impl: map[string]Connector{}}
}

func (r *Registry) Register(c Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.impl[c.Name()] = c
}

func (r *Registry) Health(ctx context.Context) []Health {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Health, 0, len(r.impl))
	for _, c := range r.impl {
		out = append(out, c.Health(ctx))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Connector < out[j].Connector })
	return out
}
