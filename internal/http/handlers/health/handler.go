package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/laoitdev/random-ml-team/internal/config"
	"github.com/laoitdev/random-ml-team/internal/random"
)

// Status represents the health state of a component.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// Check captures the outcome of an individual health probe.
type Check struct {
	Name    string                 `json:"name"`
	Status  Status                 `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// Metadata provides contextual telemetry about the service.
type Metadata struct {
	RequestCount uint64    `json:"requestCount"`
	StartedAt    time.Time `json:"startedAt"`
	Uptime       string    `json:"uptime"`
}

// Response represents the payload returned by detailed health endpoints.
type Response struct {
	Status    Status    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Checks    []Check   `json:"checks,omitempty"`
	Metadata  Metadata  `json:"metadata"`
}

// TeamGenerator exposes the subset of generator behaviour required for health checks.
type TeamGenerator interface {
	Generate() (random.Team, error)
}

// Handler exposes health endpoints with varying depth of diagnostics.
type Handler struct {
	cfg       *config.Config
	generator TeamGenerator
	metrics   *Metrics
	version   string
}

// NewHandler wires a health handler with configuration, generator and metrics dependencies.
func NewHandler(cfg *config.Config, generator TeamGenerator, metrics *Metrics, version string) *Handler {
	return &Handler{
		cfg:       cfg,
		generator: generator,
		metrics:   metrics,
		version:   version,
	}
}

// Healthz godoc
// @Summary Legacy health probe
// @Tags health
// @Produce plain
// @Success 200 {string} string "ok"
// @Router /healthz [get]
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Live godoc
// @Summary Liveness probe
// @Tags health
// @Produce json
// @Success 200 {object} Response
// @Router /health/live [get]
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	resp := h.baseResponse(StatusHealthy, nil)
	h.writeJSON(w, http.StatusOK, resp)
}

// Ready godoc
// @Summary Readiness probe
// @Tags health
// @Produce json
// @Success 200 {object} Response
// @Failure 429 {object} Response
// @Failure 503 {object} Response
// @Router /health/ready [get]
func (h *Handler) Ready(w http.ResponseWriter, _ *http.Request) {
	checks := h.coreChecks()
	h.respond(w, checks)
}

// Status godoc
// @Summary Detailed health status
// @Tags health
// @Produce json
// @Success 200 {object} Response
// @Failure 429 {object} Response
// @Failure 503 {object} Response
// @Router /health/status [get]
func (h *Handler) Status(w http.ResponseWriter, _ *http.Request) {
	checks := append(h.coreChecks(), h.runtimeCheck())
	h.respond(w, checks)
}

// Root godoc
// @Summary Detailed health status
// @Tags health
// @Produce json
// @Success 200 {object} Response
// @Failure 429 {object} Response
// @Failure 503 {object} Response
// @Router /health [get]
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	checks := append(h.coreChecks(), h.runtimeCheck())
	h.respond(w, checks)
}

func (h *Handler) baseResponse(status Status, checks []Check) Response {
	startedAt := time.Time{}
	requestCount := uint64(0)
	uptime := ""

	if h.metrics != nil {
		startedAt = h.metrics.StartedAt().UTC()
		uptime = h.metrics.Uptime().String()
		requestCount = h.metrics.RequestCount()
	}

	if uptime == "" {
		uptime = "0s"
	}

	return Response{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Version:   h.version,
		Checks:    checks,
		Metadata: Metadata{
			RequestCount: requestCount,
			StartedAt:    startedAt,
			Uptime:       uptime,
		},
	}
}

func (h *Handler) coreChecks() []Check {
	return []Check{
		h.configurationCheck(),
		h.generatorCheck(),
	}
}

func (h *Handler) configurationCheck() Check {
	check := Check{
		Name:   "configuration",
		Status: StatusHealthy,
		Details: map[string]interface{}{
			"compositionSize": len(h.cfg.Team.Composition),
			"allowDuplicates": h.cfg.Team.AllowDuplicates,
			"rolesConfigured": len(h.cfg.Team.Heroes),
		},
	}

	if len(h.cfg.Team.Composition) == 0 {
		check.Status = StatusUnhealthy
		check.Error = "team composition must not be empty"
		return check
	}

	for _, role := range h.cfg.Team.Composition {
		if len(h.cfg.Team.Heroes[role]) == 0 {
			check.Status = StatusUnhealthy
			check.Error = fmt.Sprintf("role %s has no heroes configured", role)
			return check
		}
	}

	return check
}

func (h *Handler) generatorCheck() Check {
	check := Check{
		Name:   "teamGenerator",
		Status: StatusHealthy,
	}

	if h.generator == nil {
		check.Status = StatusUnhealthy
		check.Error = "generator unavailable"
		return check
	}

	if _, err := h.generator.Generate(); err != nil {
		check.Status = StatusUnhealthy
		check.Error = err.Error()
	}

	return check
}

func (h *Handler) runtimeCheck() Check {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return Check{
		Name:   "runtime",
		Status: StatusHealthy,
		Details: map[string]interface{}{
			"allocBytes":      mem.Alloc,
			"totalAllocBytes": mem.TotalAlloc,
			"sysBytes":        mem.Sys,
			"numGC":           mem.NumGC,
			"numGoroutine":    runtime.NumGoroutine(),
			"nextGCBytes":     mem.NextGC,
		},
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respond(w http.ResponseWriter, checks []Check) {
	status := aggregateStatus(checks)
	resp := h.baseResponse(status, checks)
	code := statusCode(status)
	h.writeJSON(w, code, resp)
}

func aggregateStatus(checks []Check) Status {
	status := StatusHealthy
	for _, check := range checks {
		switch check.Status {
		case StatusUnhealthy:
			return StatusUnhealthy
		case StatusDegraded:
			status = StatusDegraded
		}
	}
	return status
}

func statusCode(status Status) int {
	switch status {
	case StatusHealthy:
		return http.StatusOK
	case StatusDegraded:
		return http.StatusTooManyRequests
	default:
		return http.StatusServiceUnavailable
	}
}
