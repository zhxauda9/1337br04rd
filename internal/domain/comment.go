package domain

import (
	"time"
)

type Comment struct {
	ID               int       `json:"id"`
	PostID           int       `json:"post_id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	AuthorID         string    `json:"author_id"`
	AuthorName       string    `json:"author_name"`
	ImageURL         string    `json:"image_url"`
	ReplyToCommentID *int      `json:"reply_to_comment_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ArchivedComment struct {
	ID               int       `json:"id"`
	PostID           int       `json:"post_id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	AuthorID         string    `json:"author_id"`
	AuthorName       string    `json:"author_name"`
	ImageURL         string    `json:"image_url"`
	CreatedAt        time.Time `json:"created_at"`
	ReplyToCommentID *int      `json:"reply_to_comment_id"`
}

func NewComment(postID int, authorID, authorName, title, content, imageURL string, replyToID *int) (*Comment, error) {
	if len(content) == 0 {
		return nil, ErrEmptyComment
	}
	return &Comment{
		PostID:           postID,
		AuthorID:         authorID,
		AuthorName:       authorName,
		Title:            title,
		Content:          content,
		ImageURL:         imageURL,
		ReplyToCommentID: replyToID,
		CreatedAt:        time.Now(),
	}, nil
}
