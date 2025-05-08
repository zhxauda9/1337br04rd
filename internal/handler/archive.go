package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"1337b04rd/internal/service"
	"1337b04rd/pkg/middleware"
)

type ArchiveHandler struct {
	service *service.ArchiveService
	logger  *slog.Logger
}

func NewArchiveHandler(service *service.ArchiveService, logger *slog.Logger) *ArchiveHandler {
	return &ArchiveHandler{service: service, logger: logger}
}

func (h *ArchiveHandler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetSession(r.Context())
	if err != nil {
		h.logger.Error("Unauthorized: session not found", "method", r.Method)
		http.Error(w, "Unauthorized: session not found", http.StatusUnauthorized)
		return
	}
	sessionID := session.ID

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		h.logger.Error("Invalid URL format", "method", r.Method)
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.logger.Error("Invalid post ID", "method", r.Method)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.service.ArchivePostByID(r.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to archive post", "error", err, "method", r.Method)
		http.Error(w, "Failed to archive post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully archived post", "method", r.Method)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post archived successfully by session " + sessionID})
}

func (h *ArchiveHandler) ArchiveExpired(w http.ResponseWriter, r *http.Request) {
	err := h.service.ArchiveExpiredPosts(r.Context())
	if err != nil {
		h.logger.Error("Failed to archive expired posts", "error", err, "method", r.Method)
		http.Error(w, "Failed to archive expired posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully archived post", "method", r.Method)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expired posts archived"})
}

func (h *ArchiveHandler) GetAllArchivedPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllArchivedPosts(r.Context())
	if err != nil {
		h.logger.Error("Failed to fetch archived posts", "error", err, "method", r.Method)
		http.Error(w, "Failed to fetch archived posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully got all archived post", "method", r.Method)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
