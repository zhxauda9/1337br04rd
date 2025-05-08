package mocks

import (
	"context"
	"errors"
	"time"

	"1337b04rd/internal/domain"
)

type mockPostRepository struct {
	posts map[int]*domain.Post
}

func NewMockPostRepository() *mockPostRepository {
	return &mockPostRepository{posts: make(map[int]*domain.Post)}
}

func (m *mockPostRepository) Save(ctx context.Context, post *domain.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	post, ok := m.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *mockPostRepository) FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Post, error) {
	var result []*domain.Post
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPostRepository) Delete(ctx context.Context, id int) error {
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepository) ArchiveExpiredPosts(ctx context.Context) error {
	return nil
}

func (m *mockPostRepository) FindExpired(ctx context.Context) ([]*domain.Post, error) {
	var expired []*domain.Post
	now := time.Now()
	for _, p := range m.posts {
		if p.ExpiresAt.Before(now) {
			expired = append(expired, p)
		}
	}
	return expired, nil
}

func (m *mockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) UpdateExpiration(ctx context.Context, postID int, newExp time.Time) error {
	if post, ok := m.posts[postID]; ok {
		post.ExpiresAt = newExp
		return nil
	}
	return errors.New("post not found")
}

func (m *mockPostRepository) UpdatePostAuthorName(ctx context.Context, postID int, authorName string) error {
	if post, ok := m.posts[postID]; ok {
		post.AuthorName = authorName
		return nil
	}
	return errors.New("post not found")
}

func (m *mockPostRepository) UpdateCommentsAuthorName(ctx context.Context, postID int, authorName string) error {
	// Only simulates update; real implementation would need commentRepo access
	if post, ok := m.posts[postID]; ok {
		post.AuthorName = authorName + " (comments updated too)"
		return nil
	}
	return errors.New("post not found")
}

type mockArchiveRepo struct {
	SaveFunc    func(ctx context.Context, post *domain.Post) error
	FindAllFunc func(ctx context.Context) ([]*domain.Archive, error)
}

func (m *mockArchiveRepo) Save(ctx context.Context, post *domain.Post) error {
	return m.SaveFunc(ctx, post)
}

func (m *mockArchiveRepo) FindAll(ctx context.Context) ([]*domain.Archive, error) {
	return m.FindAllFunc(ctx)
}
