package ports

import (
	"context"
	"time"

	"1337b04rd/internal/domain"
)

type PostRepository interface {
	Save(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id int) (*domain.Post, error)
	FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Post, error)
	Delete(ctx context.Context, id int) error
	ArchiveExpiredPosts(context.Context) error
	FindExpired(ctx context.Context) ([]*domain.Post, error)
	UpdatePostAuthorName(ctx context.Context, postID int, authorName string) error
	UpdateCommentsAuthorName(ctx context.Context, postID int, authorName string) error
	Update(ctx context.Context, comment *domain.Post) error
	UpdateExpiration(ctx context.Context, postID int, newExpiration time.Time) error
}

// want (ports.PostService, *storage.MinioClient)

type PostService interface {
	CreatePost(ctx context.Context, authorID string, title, content string) (*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
	DeleteExpiredPost(ctx context.Context, delay time.Duration) error
	// UpdatePost(ctx context.Context, id int, title, content string) error
	GetAllPosts(ctx context.Context) ([]*domain.Post, error)
	ArchiveExpiredPosts(ctx context.Context) error
}
