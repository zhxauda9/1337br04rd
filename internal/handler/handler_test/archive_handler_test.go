package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"1337b04rd/internal/domain"
)

// Mock ArchiveService
type MockArchiveService struct{}

func (s *MockArchiveService) ArchivePost(ctx interface{}, postID string) (*domain.Post, error) {
	// Simulate archiving a post by returning a mock post
	return &domain.Post{
		ID:      1,
		Title:   "Archived Post",
		Content: "Content of archived post",
	}, nil
}

func (s *MockArchiveService) GetArchivedPosts(ctx interface{}) ([]*domain.Post, error) {
	// Return mock archived posts
	return []*domain.Post{
		{ID: 1, Title: "Archived Post 1", Content: "Content 1"},
		{ID: 2, Title: "Archived Post 2", Content: "Content 2"},
	}, nil
}

// Mock ArchiveHandler
type MockArchiveHandler struct {
	service *MockArchiveService
}

func NewArchiveHandler(service *MockArchiveService) *MockArchiveHandler {
	// Return an instance of the handler
	return &MockArchiveHandler{service: service}
}

func (h *MockArchiveHandler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	// Simulate archiving the post and responding with the result
	postID := r.URL.Path[len("/archive/"):]

	// Call the ArchivePost method from the service
	post, err := h.service.ArchivePost(nil, postID)
	if err != nil {
		http.Error(w, "Failed to archive post", http.StatusInternalServerError)
		return
	}

	// Encode the response to JSON
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Test Archive Post
func TestArchivePost(t *testing.T) {
	// Create a mock service and handler
	archiveService := &MockArchiveService{}
	handler := NewArchiveHandler(archiveService)

	// Prepare the request (archiving a post with ID "test-post-id")
	req, err := http.NewRequest("POST", "/archive/test-post-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the response recorder
	rr := &TestResponseRecorder{Headers: make(http.Header)}

	// Call the handler
	handler.ArchivePost(rr, req)

	// Assert status code
	if rr.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.StatusCode)
	}

	// Assert response body
	var post domain.Post
	if err := json.NewDecoder(bytes.NewReader(rr.Body)).Decode(&post); err != nil {
		t.Fatal("Failed to decode response", err)
	}

	if post.Title != "Archived Post" {
		t.Errorf("Expected post title 'Archived Post', got '%s'", post.Title)
	}

	if post.Content != "Content of archived post" {
		t.Errorf("Expected post content 'Content of archived post', got '%s'", post.Content)
	}
}

type TestResponseRecorder struct {
	Headers    http.Header
	Body       []byte
	StatusCode int
}

func (rr *TestResponseRecorder) Header() http.Header {
	return rr.Headers
}

func (rr *TestResponseRecorder) Write(b []byte) (int, error) {
	rr.Body = append(rr.Body, b...)
	return len(b), nil
}

func (rr *TestResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
}
