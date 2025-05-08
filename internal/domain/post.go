package domain

import (
	"time"
)

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ImageURL   string    `json:"image_url"`
	Comments   []Comment `json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func NewPost(title, content string, authorID string, comments []Comment) (*Post, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	if len(title) == 0 {
		return nil, ErrTitleEmpty
	}

	now := time.Now()

	ttl := 10 * time.Minute // Default is 10 minutes

	return &Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		ImageURL:  "",
		Comments:  comments,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}, nil
}
