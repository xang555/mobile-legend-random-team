package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/laoitdev/random-ml-team/internal/random"
)

// TeamGenerator captures the behaviour required from a team generator.
type TeamGenerator interface {
	Generate() (random.Team, error)
}

// TeamHandler serves endpoints for team generation.
type TeamHandler struct {
	generator TeamGenerator
	logger    *zap.Logger
}

// NewTeamHandler creates a new TeamHandler instance.
func NewTeamHandler(generator TeamGenerator, logger *zap.Logger) *TeamHandler {
	return &TeamHandler{generator: generator, logger: logger}
}

// Register mounts the handler routes onto the provided router.
func (h *TeamHandler) Register(r chiRouter) {
	r.Get("/team/random", h.randomTeam)
}

// ErrorResponse represents a standard error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// randomTeam godoc
// @Summary Generate random team
// @Description Returns a randomly generated team based on configuration
// @Tags team
// @Produce json
// @Success 200 {object} random.Team
// @Failure 500 {object} handlers.ErrorResponse
// @Router /team/random [get]
func (h *TeamHandler) randomTeam(w http.ResponseWriter, r *http.Request) {
	team, err := h.generator.Generate()
	if err != nil {
		h.logger.Error("failed to generate team", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "could not generate team")
		return
	}

	writeJSON(w, http.StatusOK, team)
}

// chiRouter is satisfied by chi.Router; abstracted for easier testing.
type chiRouter interface {
	Get(pattern string, handlerFn http.HandlerFunc)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Response writer should not be nil; best effort logging is necessary.
		// We cannot log here since logger is not available.
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
