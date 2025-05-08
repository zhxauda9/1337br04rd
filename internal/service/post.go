package service

import (
	"context"
	"fmt"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/repository"
)

type PostService interface {
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
	DeleteExpiredPost(ctx context.Context) error
	GetAllPosts(ctx context.Context) ([]*domain.Post, error)
	ArchiveExpiredPosts(ctx context.Context) error
}

type postService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) CreatePost(ctx context.Context, post *domain.Post) error {
	// Validate the post data
	if len(post.Title) == 0 || len(post.Content) == 0 {
		return fmt.Errorf("title and content cannot be empty")
	}

	// Use repository to save post
	if err := s.postRepo.Save(ctx, post); err != nil {
		return fmt.Errorf("unable to save post: %w", err)
	}

	return nil
}

func (s *postService) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	return s.postRepo.FindByID(ctx, id)
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	return s.postRepo.Delete(ctx, id)
}

func (s *postService) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	// Retrieve all posts
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve posts: %w", err)
	}

	for _, post := range posts {
		err := s.postRepo.UpdateAuthorNameForPostAndComments(ctx, post.ID, post.AuthorName)
		if err != nil {
			return nil, fmt.Errorf("unable to update author name for post %d: %w", post.ID, err)
		}
	}
	return posts, nil
}

func (s *postService) ArchiveExpiredPosts(ctx context.Context) error {
	return s.postRepo.ArchiveExpiredPosts(ctx)
}

func (s *postService) DeleteExpiredPost(ctx context.Context) error {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch posts: %w", err)
	}

	currentTime := time.Now()

	for _, post := range posts {
		if len(post.Comments) == 0 && currentTime.Sub(post.CreatedAt) > 10*time.Minute {
			err := s.postRepo.Delete(ctx, post.ID)
			if err != nil {
				return fmt.Errorf("unable to delete post without comments: %w", err)
			}
			continue
		}

		if len(post.Comments) > 0 {
			lastComment := post.Comments[len(post.Comments)-1]
			if currentTime.Sub(lastComment.CreatedAt) > 15*time.Minute {
				err := s.postRepo.Delete(ctx, post.ID)
				if err != nil {
					return fmt.Errorf("unable to delete post with comments: %w", err)
				}
			}
		}
	}

	return nil
}
