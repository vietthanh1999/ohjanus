package service

import (
	"sort"
	"sync"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// ConnRegistry owns connection pools and their metadata behind a RWMutex so
// the Admin API can register connections at runtime while query paths read
// them concurrently. All services and the admin server share one instance.
type ConnRegistry struct {
	mu    sync.RWMutex
	pools map[string]out.Pool
	metas map[string]domain.ConnectionMeta
}

// NewConnRegistry wraps startup pools/metas. Nil maps become empty ones.
func NewConnRegistry(pools map[string]out.Pool, metas map[string]domain.ConnectionMeta) *ConnRegistry {
	if pools == nil {
		pools = map[string]out.Pool{}
	}
	if metas == nil {
		metas = map[string]domain.ConnectionMeta{}
	}
	return &ConnRegistry{pools: pools, metas: metas}
}

// Pool returns the pool for a connection alias.
func (r *ConnRegistry) Pool(name string) (out.Pool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pools[name]
	return p, ok
}

// Meta returns the metadata for a connection alias.
func (r *ConnRegistry) Meta(name string) (domain.ConnectionMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.metas[name]
	return m, ok
}

// Has reports whether a connection alias is registered.
func (r *ConnRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.metas[name]
	return ok
}

// Metas returns a name-sorted snapshot of all registered metadata.
func (r *ConnRegistry) Metas() []domain.ConnectionMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.ConnectionMeta, 0, len(r.metas))
	for _, m := range r.metas {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Connection.Name < out[j].Connection.Name })
	return out
}

// Pools returns a snapshot of the registered pools (for health checks).
func (r *ConnRegistry) Pools() map[string]out.Pool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make(map[string]out.Pool, len(r.pools))
	for n, p := range r.pools {
		cp[n] = p
	}
	return cp
}

// Attach registers (or replaces) a connection at runtime.
func (r *ConnRegistry) Attach(name string, pool out.Pool, meta domain.ConnectionMeta) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pools[name] = pool
	r.metas[name] = meta
}
