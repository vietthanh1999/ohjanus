package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Store is an in-memory TokenStore for single-instance deployments.
// Multi-instance deployments replace it with Redis (same port).
type Store struct {
	mu     sync.Mutex
	tokens map[string]*domain.PreviewToken
	subs   map[chan domain.TokenEvent]struct{}
	clock  out.Clock
	ttl    time.Duration
}

var _ out.TokenStore = (*Store)(nil)

// New builds a store with the given preview-token TTL (0 = no expiry).
func New(clock out.Clock, ttl time.Duration) *Store {
	return &Store{
		tokens: map[string]*domain.PreviewToken{},
		subs:   map[chan domain.TokenEvent]struct{}{},
		clock:  clock,
		ttl:    ttl,
	}
}

// Create stores a new pending token.
func (s *Store) Create(_ context.Context, t *domain.PreviewToken) error {
	if t == nil || t.ID == "" {
		return fmt.Errorf("token id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, dup := s.tokens[t.ID]; dup {
		return fmt.Errorf("token %q already exists", t.ID)
	}
	cp := *t
	cp.State = domain.TokenPending
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = s.clock.Now()
	}
	if s.ttl > 0 {
		cp.ExpiresAt = cp.CreatedAt.Add(s.ttl)
	}
	s.tokens[cp.ID] = &cp
	s.broadcast(domain.TokenEvent{Type: domain.TokenEventCreated, Token: copyOf(&cp)})
	return nil
}

// Get returns a copy, applying lazy expiry.
func (s *Store) Get(_ context.Context, id string) (*domain.PreviewToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[id]
	if !ok {
		return nil, domain.NewError(domain.CodeTokenMismatch, "unknown preview token")
	}
	s.expireLocked(t)
	return copyOf(t), nil
}

// Approve moves a pending token to approved (idempotent on approved).
func (s *Store) Approve(_ context.Context, id, decidedBy, reason string) (*domain.PreviewToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[id]
	if !ok {
		return nil, domain.NewError(domain.CodeTokenMismatch, "unknown preview token")
	}
	if s.expireLocked(t) {
		return nil, domain.NewError(domain.CodeTokenExpired, "preview token expired")
	}
	switch t.State {
	case domain.TokenApproved:
		return copyOf(t), nil
	case domain.TokenPending:
		t.State = domain.TokenApproved
		t.ApprovedBy = decidedBy
		t.DecidedReason = reason
		t.DecidedAt = s.clock.Now()
		s.broadcast(domain.TokenEvent{Type: domain.TokenEventApproved, Token: copyOf(t)})
		return copyOf(t), nil
	case domain.TokenUsed:
		return nil, domain.NewError(domain.CodeTokenAlreadyUsed, "preview token already used")
	default:
		return nil, domain.NewError(domain.CodeTokenMismatch, fmt.Sprintf("cannot approve token in state %q", t.State))
	}
}

// Reject moves a pending token to rejected. A reason is required.
func (s *Store) Reject(_ context.Context, id, decidedBy, reason string) (*domain.PreviewToken, error) {
	if reason == "" {
		return nil, fmt.Errorf("rejection requires a reason")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[id]
	if !ok {
		return nil, domain.NewError(domain.CodeTokenMismatch, "unknown preview token")
	}
	if s.expireLocked(t) {
		return nil, domain.NewError(domain.CodeTokenExpired, "preview token expired")
	}
	if t.State != domain.TokenPending {
		return nil, domain.NewError(domain.CodeTokenMismatch, fmt.Sprintf("cannot reject token in state %q", t.State))
	}
	t.State = domain.TokenRejected
	t.ApprovedBy = decidedBy
	t.DecidedReason = reason
	t.DecidedAt = s.clock.Now()
	s.broadcast(domain.TokenEvent{Type: domain.TokenEventRejected, Token: copyOf(t)})
	return copyOf(t), nil
}

// MarkUsed consumes an approved token (single-use).
func (s *Store) MarkUsed(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tokens[id]
	if !ok {
		return domain.NewError(domain.CodeTokenMismatch, "unknown preview token")
	}
	if s.expireLocked(t) {
		return domain.NewError(domain.CodeTokenExpired, "preview token expired")
	}
	switch t.State {
	case domain.TokenUsed:
		return domain.NewError(domain.CodeTokenAlreadyUsed, "preview token already used")
	case domain.TokenApproved:
		t.State = domain.TokenUsed
		s.broadcast(domain.TokenEvent{Type: domain.TokenEventUsed, Token: copyOf(t)})
		return nil
	default:
		return domain.NewError(domain.CodeApprovalRequired, fmt.Sprintf("token in state %q is not approved", t.State))
	}
}

// List returns copies of all tokens, newest first.
func (s *Store) List(_ context.Context) ([]*domain.PreviewToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*domain.PreviewToken, 0, len(s.tokens))
	for _, t := range s.tokens {
		s.expireLocked(t)
		out = append(out, copyOf(t))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

// Watch streams lifecycle events. Call cancel to unsubscribe.
func (s *Store) Watch() (<-chan domain.TokenEvent, func()) {
	ch := make(chan domain.TokenEvent, 16)
	s.mu.Lock()
	s.subs[ch] = struct{}{}
	s.mu.Unlock()
	return ch, func() {
		s.mu.Lock()
		delete(s.subs, ch)
		s.mu.Unlock()
	}
}

// expireLocked marks a pending token expired. Callers must hold s.mu.
func (s *Store) expireLocked(t *domain.PreviewToken) bool {
	if t.State != domain.TokenPending || t.ExpiresAt.IsZero() {
		return false
	}
	if s.clock.Now().After(t.ExpiresAt) {
		t.State = domain.TokenExpired
		s.broadcast(domain.TokenEvent{Type: domain.TokenEventExpired, Token: copyOf(t)})
		return true
	}
	return false
}

func (s *Store) broadcast(e domain.TokenEvent) {
	for ch := range s.subs {
		select {
		case ch <- e:
		default:
		}
	}
}

func copyOf(t *domain.PreviewToken) *domain.PreviewToken {
	cp := *t
	if t.Params != nil {
		cp.Params = append([]any(nil), t.Params...)
	}
	if t.Warnings != nil {
		cp.Warnings = append([]string(nil), t.Warnings...)
	}
	return &cp
}
