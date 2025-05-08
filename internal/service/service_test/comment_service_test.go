package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/service"
)

type mockCommentRepo struct {
	comments map[int]*domain.Comment
}

func (m *mockCommentRepo) Save(ctx context.Context, c *domain.Comment) error {
	m.comments[c.ID] = c
	return nil
}

func (m *mockCommentRepo) FindByID(ctx context.Context, id int) (*domain.Comment, error) {
	comment, exists := m.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (m *mockCommentRepo) FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Comment, error) {
	return nil, nil
}

func (m *mockCommentRepo) FindCommentOfPost(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return nil, nil
}

func (m *mockCommentRepo) FindRepliesToComment(ctx context.Context, commentID int) ([]*domain.Comment, error) {
	return nil, nil
}

func (m *mockCommentRepo) Update(ctx context.Context, c *domain.Comment) error {
	return nil
}

type mockPostRepo struct {
	posts map[int]*domain.Post
}

func (m *mockPostRepo) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *mockPostRepo) UpdateExpiration(ctx context.Context, postID int, newExpiration time.Time) error {
	post, exists := m.posts[postID]
	if !exists {
		return errors.New("post not found")
	}
	post.ExpiresAt = newExpiration
	return nil
}

// unused methods
func (m *mockPostRepo) Save(ctx context.Context, p *domain.Post) error { return nil }

func (m *mockPostRepo) FindByAuthorID(ctx context.Context, id string) ([]*domain.Post, error) {
	return nil, nil
}
func (m *mockPostRepo) Delete(ctx context.Context, id int) error                { return nil }
func (m *mockPostRepo) ArchiveExpiredPosts(ctx context.Context) error           { return nil }
func (m *mockPostRepo) FindExpired(ctx context.Context) ([]*domain.Post, error) { return nil, nil }
func (m *mockPostRepo) UpdatePostAuthorName(ctx context.Context, id int, name string) error {
	return nil
}

func (m *mockPostRepo) UpdateCommentsAuthorName(ctx context.Context, id int, name string) error {
	return nil
}
func (m *mockPostRepo) Update(ctx context.Context, p *domain.Post) error { return nil }

func TestCommentService_CreateComment_Success(t *testing.T) {
	post := &domain.Post{ID: 1, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(10 * time.Minute)}
	postRepo := &mockPostRepo{posts: map[int]*domain.Post{1: post}}
	commentRepo := &mockCommentRepo{comments: make(map[int]*domain.Comment)}
	svc := service.NewCommentService(commentRepo, postRepo)

	comment, err := svc.CreateComment(context.Background(), 1, "user1", "Author", "Title", "Content", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if comment.PostID != 1 || comment.AuthorID != "user1" {
		t.Errorf("unexpected comment result: %+v", comment)
	}

	if post.ExpiresAt.Sub(post.CreatedAt) < 15*time.Minute {
		t.Errorf("post expiration not extended properly")
	}
}

func TestCommentService_CreateComment_PostNotFound(t *testing.T) {
	postRepo := &mockPostRepo{posts: make(map[int]*domain.Post)}
	commentRepo := &mockCommentRepo{comments: make(map[int]*domain.Comment)}
	svc := service.NewCommentService(commentRepo, postRepo)

	_, err := svc.CreateComment(context.Background(), 999, "user1", "Author", "Title", "Content", "", nil)
	if err == nil {
		t.Fatal("expected error when post does not exist")
	}
}

func TestCommentService_CreateComment_ReplyNotFound(t *testing.T) {
	post := &domain.Post{ID: 1, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(10 * time.Minute)}
	postRepo := &mockPostRepo{posts: map[int]*domain.Post{1: post}}
	commentRepo := &mockCommentRepo{comments: make(map[int]*domain.Comment)}
	svc := service.NewCommentService(commentRepo, postRepo)

	replyID := 999
	_, err := svc.CreateComment(context.Background(), 1, "user1", "Author", "Title", "Content", "", &replyID)
	if err == nil {
		t.Fatal("expected error when reply-to comment does not exist")
	}
}
