package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/laoitdev/random-ml-team/internal/config"
	"github.com/laoitdev/random-ml-team/internal/random"
)

type stubGenerator struct {
	team random.Team
	err  error
}

func (s stubGenerator) Generate() (random.Team, error) {
	return s.team, s.err
}

func TestReadyHealthy(t *testing.T) {
	cfg := &config.Config{
		Team: config.Team{
			Composition:     []string{"tank"},
			AllowDuplicates: false,
			Heroes: map[string][]string{
				"tank": {"hero"},
			},
		},
	}

	metrics := NewMetrics(time.Now().Add(-time.Minute))
	handler := NewHandler(cfg, stubGenerator{team: random.Team{}}, metrics, "test")

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	resp := httptest.NewRecorder()

	handler.Ready(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", resp.Code)
	}

	var payload Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Status != StatusHealthy {
		t.Fatalf("expected overall status healthy got %s", payload.Status)
	}

	if payload.Metadata.RequestCount != 0 {
		t.Errorf("expected request count 0 got %d", payload.Metadata.RequestCount)
	}
}

func TestReadyUnhealthyWhenGeneratorFails(t *testing.T) {
	cfg := &config.Config{
		Team: config.Team{
			Composition:     []string{"tank"},
			AllowDuplicates: false,
			Heroes: map[string][]string{
				"tank": {"hero"},
			},
		},
	}

	handler := NewHandler(cfg, stubGenerator{err: errors.New("boom")}, NewMetrics(time.Now()), "test")
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)

	handler.Ready(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503 got %d", resp.Code)
	}

	var payload Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy got %s", payload.Status)
	}

	found := false
	for _, check := range payload.Checks {
		if check.Name == "teamGenerator" {
			found = true
			if check.Status != StatusUnhealthy {
				t.Fatalf("expected generator check unhealthy got %s", check.Status)
			}
		}
	}

	if !found {
		t.Fatal("expected generator check present")
	}
}

func TestRootDetailedStatus(t *testing.T) {
	cfg := &config.Config{
		Team: config.Team{
			Composition:     []string{"tank"},
			AllowDuplicates: false,
			Heroes: map[string][]string{
				"tank": {"hero"},
			},
		},
	}

	handler := NewHandler(cfg, stubGenerator{team: random.Team{}}, NewMetrics(time.Now()), "v1")
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.Root(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", resp.Code)
	}

	var payload Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(payload.Checks) == 0 {
		t.Fatal("expected checks in detailed response")
	}
}

func TestMetricsTracksRequests(t *testing.T) {
	start := time.Now().Add(-5 * time.Minute)
	metrics := NewMetrics(start)
	metrics.IncrementRequests()
	metrics.IncrementRequests()

	if metrics.RequestCount() != 2 {
		t.Fatalf("expected 2 requests got %d", metrics.RequestCount())
	}

	if metrics.StartedAt() != start {
		t.Fatalf("expected start time %v got %v", start, metrics.StartedAt())
	}

	if metrics.Uptime() <= 0 {
		t.Fatal("expected positive uptime")
	}
}
