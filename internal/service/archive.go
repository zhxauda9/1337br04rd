package service

import (
	"context"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/repository"
)

type ArchiveService struct {
	archiveRepo *repository.ArchiveRepository
	postRepo    *repository.PostRepository
}

func NewArchiveService(archiveRepo *repository.ArchiveRepository, postRepo *repository.PostRepository) *ArchiveService {
	return &ArchiveService{archiveRepo: archiveRepo, postRepo: postRepo}
}

func (s *ArchiveService) ArchivePostByID(ctx context.Context, postID int) error {
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	return s.archiveRepo.Save(ctx, post)
}

func (s *ArchiveService) ArchiveExpiredPosts(ctx context.Context) error {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return err
	}
	for _, post := range posts {
		if post.ExpiresAt.Before(time.Now()) {
			if err := s.archiveRepo.Save(ctx, post); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ArchiveService) GetAllArchivedPosts(ctx context.Context) ([]*domain.Archive, error) {
	return s.archiveRepo.FindAll(ctx)
}

func (s *ArchiveService) GetArchivedPostByID(ctx context.Context, postID int) *domain.Post {
	post, _ := s.archiveRepo.FindByID(ctx, postID)
	return post
}
