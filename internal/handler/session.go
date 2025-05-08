package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"1337b04rd/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
	logger         *slog.Logger
}

func NewSessionHandler(sessionService *service.SessionService, logger *slog.Logger) *SessionHandler {
	return &SessionHandler{sessionService: sessionService, logger: logger}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body", "method", r.Method)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	session, err := h.sessionService.CreateSession(r.Context(), req.Name)
	if err != nil {
		h.logger.Error("Failed to create session", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_name",
		Value:    session.Name,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	h.logger.Info("Successfully created session", "method", r.Method)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	err := h.sessionService.DeleteSession(w, r)
	if err != nil {
		h.logger.Error("Failed to delete session", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Failed to delete session: %v", err), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully deleted session", "method", r.Method)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.logger.Error("No session cookie found", "error", err, "method", r.Method)
		http.Error(w, "No session cookie found", http.StatusUnauthorized)
		return
	}

	session, err := h.sessionService.GetSession(r.Context(), cookie.Value)
	if err != nil {
		h.logger.Error("Failed to get session", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Failed to get session: %v", err), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully got session", "method", r.Method)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (h *SessionHandler) GetAllSession(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetAllSessions(r.Context())
	if err != nil {
		h.logger.Error("Failed to get sessions", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Failed to get sessions: %v", err), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully got all sessions", "method", r.Method)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.logger.Error("No session cookie found", "error", err, "method", r.Method)
		http.Error(w, "No session cookie found", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"Name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("No session cookie found", "error", err, "method", r.Method)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Name == "" {
		h.logger.Error("Name cannot be empty", "method", r.Method)
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.UpdateSession(cookie.Value, req.Name)
	if err != nil {
		h.logger.Error("Failed to update session", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Failed to update session: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully updated session", "method", r.Method)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}
