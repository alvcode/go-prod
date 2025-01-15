package metric

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
}

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Print("heartbeat")
	w.WriteHeader(http.StatusNoContent)
}
