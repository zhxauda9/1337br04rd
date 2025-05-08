package mocks

import (
	"context"
)

type MockMinioClient struct{}

func (m *MockMinioClient) UploadCommentImage(ctx context.Context, fileBytes []byte, filename, contentType string) (string, error) {
	return "http://localhost:9001/comments/test-image.jpg", nil
}

func (m *MockMinioClient) UploadPostImage(ctx context.Context, fileBytes []byte, filename, contentType string) (string, error) {
	return "http://localhost:9001/posts/test-image.jpg", nil
}

func (m *MockMinioClient) UploadAvatarFromURL(ctx context.Context, imageUrl string) (string, error) {
	return "http://localhost:9001/avatars/test-avatar.jpg", nil
}

func (m *MockMinioClient) SaveSession(ctx context.Context, sessionID string, sessionData []byte) error {
	return nil
}

type FakeMinioClient struct{}

func (m *FakeMinioClient) UploadPostImage(ctx context.Context, fileBytes []byte, filename, contentType string) (string, error) {
	return "http://localhost:9001/posts/test-image.jpg", nil
}

func (m *FakeMinioClient) UploadCommentImage(ctx context.Context, fileBytes []byte, filename, contentType string) (string, error) {
	return "http://localhost:9001/comments/test-image.jpg", nil
}

func (m *FakeMinioClient) UploadAvatarFromURL(ctx context.Context, imageUrl string) (string, error) {
	return "http://localhost:9001/avatars/test-avatar.jpg", nil
}

func (m *FakeMinioClient) SaveSession(ctx context.Context, sessionID string, sessionData []byte) error {
	return nil
}
