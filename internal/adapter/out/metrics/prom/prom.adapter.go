package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Metrics implements out.Metrics with Prometheus counters (§13.1).
// Each instance registers its own registry, so tests stay isolated.
type Metrics struct {
	registry *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	denials  *prometheus.CounterVec
}

var _ out.Metrics = (*Metrics)(nil)

// New builds and registers the Janus metrics.
func New() *Metrics {
	reg := prometheus.NewRegistry()
	m := &Metrics{
		registry: reg,
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "janus_requests_total",
			Help: "Total MCP tool calls by tool and status.",
		}, []string{"tool", "status"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "janus_query_duration_seconds",
			Help:    "Query duration by connection and statement type.",
			Buckets: prometheus.DefBuckets,
		}, []string{"connection", "statement_type"}),
		denials: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "janus_policy_denials_total",
			Help: "Policy denials by rule and connection.",
		}, []string{"rule", "connection"}),
	}
	reg.MustRegister(m.requests, m.duration, m.denials)
	return m
}

// IncRequests implements Metrics.
func (m *Metrics) IncRequests(tool, status string) {
	m.requests.WithLabelValues(tool, status).Inc()
}

// ObserveQueryDuration implements Metrics.
func (m *Metrics) ObserveQueryDuration(connection, statement string, seconds float64) {
	m.duration.WithLabelValues(connection, statement).Observe(seconds)
}

// IncPolicyDenials implements Metrics.
func (m *Metrics) IncPolicyDenials(rule, connection string) {
	m.denials.WithLabelValues(rule, connection).Inc()
}

// Handler serves the metrics endpoint.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
