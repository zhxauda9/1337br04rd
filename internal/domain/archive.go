package domain

import (
	"time"
)

type Archive struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ImageURL   string    `json:"image_url"`
	Comments   []Comment `json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expires_at"`
	ArchivedAt time.Time `json:"archived_at"`
}

func NewArchive(title, content string, authorID string, image_url string, createdAt time.Time, expiredAt time.Time, comments []Comment) (*Archive, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	if len(title) == 0 {
		return nil, ErrTitleEmpty
	}

	now := time.Now()

	return &Archive{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		ImageURL:   image_url,
		Comments:   comments,
		CreatedAt:  createdAt,
		ExpiredAt:  expiredAt,
		ArchivedAt: now,
	}, nil
}
