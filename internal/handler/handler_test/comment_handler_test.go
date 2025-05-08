package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"1337b04rd/internal/adapters/storage"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/handler"
	"1337b04rd/pkg/logger"
)

// ---- Mock Comment Service ----

type MockCommentService struct{}

func (s *MockCommentService) CreateComment(ctx context.Context, postID int, title, content, imageURL, sessionID, userID string, parentID *int) (*domain.Comment, error) {
	if postID == 0 {
		return nil, errors.New("post not found")
	}
	return &domain.Comment{
		ID:       1,
		PostID:   postID,
		Title:    title,
		Content:  content,
		ImageURL: imageURL,
	}, nil
}

func (s *MockCommentService) GetRepliesToComment(ctx context.Context, commentID int) ([]*domain.Comment, error) {
	return []*domain.Comment{
		{ID: 1, PostID: 1, Title: "Reply 1", Content: "Reply 1 Content"},
		{ID: 2, PostID: 1, Title: "Reply 2", Content: "Reply 2 Content"},
	}, nil
}

func (s *MockCommentService) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.Comment{ID: id, PostID: 1, Title: "Test", Content: "Content"}, nil
}

func (s *MockCommentService) GetCommentsByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return []*domain.Comment{
		{ID: 1, PostID: postID, Title: "Comment 1", Content: "Content 1"},
	}, nil
}

func (s *MockCommentService) GetRepliesByParentID(ctx context.Context, commentID int) ([]*domain.Comment, error) {
	return []*domain.Comment{
		{ID: 2, PostID: 1, Title: "Reply", Content: "Reply content"},
	}, nil
}

func (s *MockCommentService) GetAllCommentsOfPost(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return []*domain.Comment{
		{ID: 1, PostID: postID, Title: "Comment 1", Content: "Content 1"},
		{ID: 2, PostID: postID, Title: "Comment 2", Content: "Content 2"},
	}, nil
}

func (s *MockCommentService) GetComment(ctx context.Context, id int) (*domain.Comment, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.Comment{
		ID:      id,
		PostID:  1,
		Title:   "Sample Comment",
		Content: "This is a sample comment.",
	}, nil
}

// ---- Mock Storage ----

type MockStorage struct{}

func (m *MockStorage) UploadCommentImage(ctx context.Context, data []byte, filename, contentType string) (string, error) {
	return "http://mocked-storage/comment.jpg", nil
}

func (m *MockStorage) UploadFile(ctx context.Context, data []byte, filename, contentType string) (string, error) {
	return "http://mocked-storage/file.jpg", nil
}

// ---- New Mock Comment Handler ----

func NewMockCommentHandler() *handler.CommentHandler {
	log, _ := logger.SetupLogger() // or discard the log file if not needed
	service := &MockCommentService{}
	storage := &storage.MinioClient{}
	return handler.NewCommentHandler(service, storage, log)
}

// ---- Tests ----

func TestGetComment_NotFound(t *testing.T) {
	h := NewMockCommentHandler()

	req := httptest.NewRequest("GET", "/comments/99", nil)
	req.URL.Path = "/comments/99"
	rr := httptest.NewRecorder()
	h.GetComment(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestGetCommentsOfPost(t *testing.T) {
	h := NewMockCommentHandler()

	req := httptest.NewRequest("GET", "/comments/post/1", nil)
	req.URL.Path = "/comments/post/1"
	rr := httptest.NewRecorder()
	h.GetCommentsOfPost(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
