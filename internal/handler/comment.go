package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"1337b04rd/internal/adapters/storage"
	"1337b04rd/internal/ports"
	"1337b04rd/pkg/middleware"
)

type CommentHandler struct {
	commentService ports.CommentService
	storage        *storage.MinioClient
	logger         *slog.Logger
}

func NewCommentHandler(commentService ports.CommentService, storage *storage.MinioClient, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{commentService: commentService, storage: storage, logger: logger}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("Could not parse multipart form", "method", r.Method)
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		h.logger.Error("Invalid post_id", "method", r.Method)
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	replyToStr := r.FormValue("reply_to_id")

	var replyToID *int
	if replyToStr != "" {
		if id, err := strconv.Atoi(replyToStr); err == nil {
			replyToID = &id
		}
	}

	session, err := middleware.GetSession(r.Context())
	if err != nil {
		h.logger.Error("Unauthorized", "method", r.Method)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var imageURL string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, file)
		if err != nil {
			h.logger.Error("Failed to read image", "method", r.Method)
			http.Error(w, "Failed to read image", http.StatusInternalServerError)
			return
		}

		imageURL, err = h.storage.UploadCommentImage(
			r.Context(),
			buf.Bytes(),
			handler.Filename,
			handler.Header.Get("Content-Type"),
		)
		if err != nil {
			h.logger.Error("Error uploading image", "error", err, "method", r.Method)
			http.Error(w, fmt.Sprintf("Error uploading image: %s", err), http.StatusInternalServerError)
			return
		}
	}

	comment, err := h.commentService.CreateComment(
		r.Context(),
		postID,
		session.ID,
		session.Name,
		title,
		content,
		imageURL,
		replyToID,
	)
	fmt.Println(comment)
	if err != nil {
		h.logger.Error("error", "error", err.Error(), "method", r.Method)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully created comment", "method", r.Method)
	http.Redirect(w, r, fmt.Sprintf("/posts/%d", postID), http.StatusSeeOther)
}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comments/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid comment ID", "method", r.Method)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.GetComment(r.Context(), id)
	if err != nil {
		h.logger.Error("Comment not found", "method", r.Method)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	h.logger.Info("Successfully got comment by id", "method", r.Method)
	writeJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) GetCommentsOfPost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comments/post/")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid post ID", "method", r.Method)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := h.commentService.GetAllCommentsOfPost(r.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to fetch comments", "method", r.Method)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully got comments of post", "method", r.Method)
	writeJSON(w, http.StatusOK, comments)
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (h *CommentHandler) GetRepliesToComment(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.logger.Error("Invalid URL", "method", r.Method)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid comment ID", "method", r.Method)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	replies, err := h.commentService.GetRepliesToComment(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to fetch replies", "method", r.Method)
		http.Error(w, "Failed to fetch replies", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully got replies to comment", "method", r.Method)

	writeJSON(w, http.StatusOK, replies)
}
