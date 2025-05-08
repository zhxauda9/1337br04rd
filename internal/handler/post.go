package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"1337b04rd/internal/adapters/storage"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"

	middleware "1337b04rd/pkg/middleware"
)

type PostHandler struct {
	postService service.PostService
	storage     *storage.MinioClient
	logger      *slog.Logger
}

func NewPostHandler(postService service.PostService, storage *storage.MinioClient, logger *slog.Logger) *PostHandler {
	return &PostHandler{postService: postService, storage: storage, logger: logger}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) == 0 || len(content) == 0 {
		http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if an image file is provided in the form
	var imageURL string
	file, _, err := r.FormFile("image")
	if err == nil {
		// If image is provided, upload it
		defer file.Close()

		session, err := middleware.GetSession(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving session: %s", err), http.StatusUnauthorized)
			return
		}

		post, err := domain.NewPost(title, content, session.ID, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating post: %s", err), http.StatusBadRequest)
			return
		}

		post.AuthorName = session.Name

		imageURL, err = h.storage.UploadImage(r.Context(), file, "post-image.jpg", "image/jpeg")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error uploading image: %s", err), http.StatusInternalServerError)
			return
		}
		post.ImageURL = imageURL
	} else if err != http.ErrMissingFile {
		// If error is something other than missing file, handle it
		http.Error(w, fmt.Sprintf("Error extracting image: %s", err), http.StatusBadRequest)
		return
	}

	session, err := middleware.GetSession(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving session: %s", err), http.StatusUnauthorized)
		return
	}

	// Create the post with or without the image URL
	post, err := domain.NewPost(title, content, session.ID, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %s", err), http.StatusBadRequest)
		return
	}

	post.AuthorName = session.Name

	// Only set ImageURL if the image was provided
	if imageURL != "" {
		post.ImageURL = imageURL
	}

	if err := h.postService.CreatePost(r.Context(), post); err != nil {
		http.Error(w, fmt.Sprintf("Error saving post: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		h.logger.Error("Error encoding response", "method", r.Method)
		http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		h.logger.Error("Missing id", "method", r.Method)
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	if idStr == "" {
		h.logger.Error("Missing id", "method", r.Method)
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid id format", "error", err, "method", r.Method)
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPost(r.Context(), intID)
	if err != nil {
		h.logger.Error("Post not found:", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Post not found: %v", err), http.StatusNotFound)
		return
	}

	if post == nil {
		h.logger.Error("Post is empty:", "method", r.Method)
		http.Error(w, "Post is empty", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully got post", "method", r.Method)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding post: %v", err), http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPosts(r.Context())
	if err != nil {
		h.logger.Error("Error retrieving posts", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Error retrieving posts: %s", err), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		h.logger.Error("No posts found", "method", r.Method)
		http.Error(w, "No posts found", http.StatusNotFound)
		return
	}

	h.logger.Info("Successfully got all posts", "method", r.Method)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.logger.Error("Error encoding posts", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Error encoding posts: %s", err), http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdateAuthorNameForAllPosts(w http.ResponseWriter, r *http.Request) {
	// Get all posts
	posts, err := h.postService.GetAllPosts(r.Context())
	if err != nil {
		h.logger.Error("Error retrieving posts", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Error retrieving posts: %s", err), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully updated author name of all posts", "method", r.Method)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.logger.Error("Error encoding posts", "error", err, "method", r.Method)
		http.Error(w, fmt.Sprintf("Error encoding posts: %s", err), http.StatusInternalServerError)
	}
}
