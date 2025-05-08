package service_test

import (
	"context"
	"testing"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/service/mocks"
)

type MockPostService struct{}

func (s *MockPostService) CreatePost(ctx context.Context, sessionID, title, content, filename, contentType string, image []byte) (*domain.Post, error) {
	return &domain.Post{
		ID:       1,
		Title:    title,
		Content:  content,
		ImageURL: "http://localhost:9001/posts/test-image.jpg",
	}, nil
}

func (s *MockPostService) GetPost(ctx context.Context, postID int) (*domain.Post, error) {
	return &domain.Post{ID: postID, Title: "Test Post", Content: "Post Content"}, nil
}

func TestUpdateCommentsAuthorName(t *testing.T) {
	repo := mocks.NewMockPostRepository()
	post := &domain.Post{
		ID:         1,
		Title:      "Hello",
		Content:    "World",
		AuthorID:   "u123",
		AuthorName: "Old Name",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}
	repo.Save(context.Background(), post)

	err := repo.UpdateCommentsAuthorName(context.Background(), 1, "New Name")
	if err != nil {
		t.Fatalf("UpdateCommentsAuthorName() failed: %v", err)
	}

	updated, _ := repo.FindByID(context.Background(), 1)
	expected := "New Name (comments updated too)"
	if updated.AuthorName != expected {
		t.Errorf("Expected author name '%s', got '%s'", expected, updated.AuthorName)
	}
}
