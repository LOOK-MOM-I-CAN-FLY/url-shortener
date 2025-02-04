package handler

import (
	"encoding/json"
	"net/http"

	"url-shortener/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Handler struct {
	service service.Shortener
	logger  *zap.Logger
}

func NewHandler(s service.Shortener, l *zap.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Decode error", zap.Error(err))
		renderError(w, http.StatusBadRequest, "invalid request")
		return
	}

	shortURL, err := h.service.CreateShortURL(r.Context(), req.URL)
	if err != nil {
		h.logger.Error("Create error", zap.Error(err))
		renderError(w, http.StatusInternalServerError, "internal error")
		return
	}

	render.JSON(w, r, map[string]string{
		"short_url": shortURL,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "shortCode")
	originalURL, err := h.service.GetOriginalURL(r.Context(), code)

	if err != nil {
		h.logger.Warn("URL not found", zap.String("code", code))
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func renderError(w http.ResponseWriter, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}
