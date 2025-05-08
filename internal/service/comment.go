package service

import (
	"context"
	"fmt"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/ports"
)

type CommentService struct {
	commentRepo ports.CommentRepository
	postRepo    ports.PostRepository // <-- Added
}

func NewCommentService(commentRepo ports.CommentRepository, postRepo ports.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, postID int, authorID, authorName, title, content, imageURL string, replyToID *int) (*domain.Comment, error) {
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, fmt.Errorf("post with ID %d not found", postID)
	}

	if replyToID != nil {
		reply, err := s.commentRepo.FindByID(ctx, *replyToID)
		if err != nil || reply == nil {
			return nil, fmt.Errorf("reply-to comment with ID %d not found", *replyToID)
		}
	}

	comment, err := domain.NewComment(postID, authorID, authorName, title, content, imageURL, replyToID)
	if err != nil {
		return nil, fmt.Errorf("create comment failed: %w", err)
	}

	if err := s.commentRepo.Save(ctx, comment); err != nil {
		return nil, err
	}

	createdPlus15 := post.CreatedAt.Add(15 * time.Minute)
	if post.ExpiresAt.Before(createdPlus15) {
		post.ExpiresAt = createdPlus15
		_ = s.postRepo.UpdateExpiration(ctx, postID, post.ExpiresAt)
	}

	return comment, nil
}

func (s *CommentService) GetComment(ctx context.Context, id int) (*domain.Comment, error) {
	return s.commentRepo.FindByID(ctx, id)
}

func (s *CommentService) GetAllCommentsOfPost(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return s.commentRepo.FindCommentOfPost(ctx, postID)
}

func (s *CommentService) GetRepliesToComment(ctx context.Context, commentID int) ([]*domain.Comment, error) {
	return s.commentRepo.FindRepliesToComment(ctx, commentID)
}
