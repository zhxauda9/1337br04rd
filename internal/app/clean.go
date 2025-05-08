package app

import (
	"context"
	"log"
	"time"

	"1337b04rd/internal/service"
)

type CleanupService struct {
	postService service.PostService
}

func NewCleanupService(postService service.PostService) *CleanupService {
	return &CleanupService{
		postService: postService,
	}
}

func (cs *CleanupService) StartCleanupTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := cs.postService.DeleteExpiredPost(context.Background()); err != nil {
				log.Println("Error deleting expired posts:", err)
			}
		}
	}()
}
