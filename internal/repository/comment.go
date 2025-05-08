package repository

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Save(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments 
			(post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.Title,
		comment.Content,
		comment.AuthorID,
		comment.AuthorName,
		comment.ImageURL,
		comment.ReplyToCommentID,
		comment.CreatedAt,
	).Scan(&comment.ID)
}

func (r *CommentRepository) FindByID(ctx context.Context, id int) (*domain.Comment, error) {
	query := `
		SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at 
		FROM comments 
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var c domain.Comment
	err := row.Scan(&c.ID, &c.PostID, &c.Title, &c.Content, &c.AuthorID, &c.AuthorName, &c.ImageURL, &c.ReplyToCommentID, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *CommentRepository) FindCommentOfPost(ctx context.Context, postID int) ([]*domain.Comment, error) {
	query := `
		SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at 
		FROM comments 
		WHERE post_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.Title, &c.Content, &c.AuthorID, &c.AuthorName, &c.ImageURL, &c.ReplyToCommentID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *CommentRepository) FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Comment, error) {
	query := `
		SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at 
		FROM comments 
		WHERE author_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.Title, &c.Content, &c.AuthorID, &c.AuthorName, &c.ImageURL, &c.ReplyToCommentID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	query := `UPDATE comments SET author_name = $1 WHERE author_id = $2`
	_, err := r.db.ExecContext(ctx, query, comment.AuthorName, comment.AuthorID)
	return err
}

func (r *CommentRepository) FindRepliesToComment(ctx context.Context, parentCommentID int) ([]*domain.Comment, error) {
	query := `
        SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at 
        FROM comments 
        WHERE reply_to_comment_id = $1
        ORDER BY created_at ASC
    `
	rows, err := r.db.QueryContext(ctx, query, parentCommentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.Title, &c.Content, &c.AuthorID, &c.AuthorName, &c.ImageURL, &c.ReplyToCommentID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		replies = append(replies, &c)
	}

	return replies, nil
}
