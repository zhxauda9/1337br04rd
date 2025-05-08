package repository

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain"
)

type ArchiveRepository struct {
	db *sql.DB
}

func NewArchiveRepository(db *sql.DB) *ArchiveRepository {
	return &ArchiveRepository{db: db}
}

func (r *ArchiveRepository) Save(ctx context.Context, post *domain.Post) error {
	// 1. Insert post into archived_posts
	insertQuery := `
        INSERT INTO archived_posts 
        (title, content, author_id, author_name, image_url, created_at, expires_at, archived_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING id
    `
	var archiveID int
	err := r.db.QueryRowContext(ctx, insertQuery,
		post.Title, post.Content, post.AuthorID, post.AuthorName, post.ImageURL,
		post.CreatedAt, post.ExpiresAt,
	).Scan(&archiveID)
	if err != nil {
		return fmt.Errorf("unable to archive post: %w", err)
	}

	// 2. Archive comments (assign new post_id pointing to archived_posts)
	for _, c := range post.Comments {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO archived_comments (
				post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, archiveID, c.Title, c.Content, c.AuthorID, c.AuthorName, c.ImageURL, c.ReplyToCommentID, c.CreatedAt)
		if err != nil {
			return fmt.Errorf("unable to archive comment: %w", err)
		}
	}

	// 3. Delete original comments
	_, err = r.db.ExecContext(ctx, "DELETE FROM comments WHERE post_id = $1", post.ID)
	if err != nil {
		return fmt.Errorf("unable to delete original comments: %w", err)
	}

	// 4. Delete original post
	_, err = r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", post.ID)
	if err != nil {
		return fmt.Errorf("unable to delete original post: %w", err)
	}

	return nil
}

func (r *ArchiveRepository) FindAll(ctx context.Context) ([]*domain.Archive, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, title, content, author_id, author_name, image_url, created_at, expires_at, archived_at
        FROM archived_posts
    `)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch archived posts: %w", err)
	}
	defer rows.Close()

	var posts []*domain.Archive
	for rows.Next() {
		var post domain.Archive
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.ImageURL,
			&post.CreatedAt,
			&post.ExpiredAt,
			&post.ArchivedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan archived post: %w", err)
		}

		// Fetch archived comments for each post
		commentRows, err := r.db.QueryContext(ctx, `
            SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at
            FROM archived_comments
            WHERE post_id = $1
        `, post.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch comments for post %d: %w", post.ID, err)
		}

		var comments []domain.Comment
		for commentRows.Next() {
			var comment domain.Comment
			var replyTo *int
			if err := commentRows.Scan(
				&comment.ID,
				&comment.PostID,
				&comment.Title,
				&comment.Content,
				&comment.AuthorID,
				&comment.AuthorName,
				&comment.ImageURL,
				&replyTo,
				&comment.CreatedAt,
			); err != nil {
				commentRows.Close()
				return nil, fmt.Errorf("unable to scan archived comment: %w", err)
			}
			comment.ReplyToCommentID = replyTo
			comments = append(comments, comment)
		}
		commentRows.Close()

		post.Comments = comments
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *ArchiveRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	query := `
        SELECT id, title, content, author_id, author_name, image_url, created_at, expires_at
        FROM archived_posts
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	var post domain.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.ImageURL, &post.CreatedAt, &post.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("archived post with id %d not found", id)
		}
		return nil, fmt.Errorf("error scanning archived post: %w", err)
	}

	// Fetch archived comments
	commentQuery := `
        SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at
        FROM archived_comments
        WHERE post_id = $1
    `
	rows, err := r.db.QueryContext(ctx, commentQuery, id)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch archived comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		var replyTo *int

		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.Title,
			&c.Content,
			&c.AuthorID,
			&c.AuthorName,
			&c.ImageURL,
			&replyTo,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan archived comment: %w", err)
		}

		c.ReplyToCommentID = replyTo
		comments = append(comments, c)
	}

	post.Comments = comments
	return &post, nil
}
