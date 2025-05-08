package ports

import (
	"context"

	"1337b04rd/internal/domain"
)

type CommentRepository interface {
	Save(ctx context.Context, comment *domain.Comment) error
	FindByID(ctx context.Context, id int) (*domain.Comment, error)
	FindCommentOfPost(ctx context.Context, postID int) ([]*domain.Comment, error)
	FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) error
	FindRepliesToComment(ctx context.Context, parentCommentID int) ([]*domain.Comment, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, postID int, authorID, authorName, title, content, imageURL string, replyToID *int) (*domain.Comment, error)
	GetComment(ctx context.Context, id int) (*domain.Comment, error)
	GetAllCommentsOfPost(ctx context.Context, postID int) ([]*domain.Comment, error)
	GetRepliesToComment(ctx context.Context, commentID int) ([]*domain.Comment, error)
}
