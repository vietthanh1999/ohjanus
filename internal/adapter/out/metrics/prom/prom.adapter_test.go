package prom

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsScrape(t *testing.T) {
	m := New()
	m.IncRequests("db_read", "success")
	m.IncRequests("db_read", "success")
	m.IncRequests("db_read", "denied")
	m.ObserveQueryDuration("analytics", "SELECT", 0.012)
	m.IncPolicyDenials("deny-ddl", "analytics")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)

	for _, want := range []string{
		`janus_requests_total{status="success",tool="db_read"} 2`,
		`janus_requests_total{status="denied",tool="db_read"} 1`,
		`janus_policy_denials_total{connection="analytics",rule="deny-ddl"} 1`,
		`janus_query_duration_seconds_count{connection="analytics",statement_type="SELECT"} 1`,
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("scrape missing %q\n%s", want, body)
		}
	}
}
