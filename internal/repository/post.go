package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"1337b04rd/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
	// If the post ID is zero, this is a new post, so we insert it
	if post.ID == 0 {
		insertPostQuery := `
            INSERT INTO posts (title, content, author_id, author_name, image_url, created_at, expires_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        `
		err := r.db.QueryRowContext(ctx, insertPostQuery, post.Title, post.Content, post.AuthorID, post.AuthorName, post.ImageURL, post.CreatedAt, post.ExpiresAt).Scan(&post.ID)
		if err != nil {
			return fmt.Errorf("unable to insert post: %w", err)
		}
	} else {
		updatePostQuery := "UPDATE posts SET title = $1, content = $2, image_url = $3 WHERE id = $4"
		_, err := r.db.ExecContext(ctx, updatePostQuery, post.Title, post.Content, post.ImageURL, post.ID)
		if err != nil {
			return fmt.Errorf("unable to update post: %w", err)
		}
	}

	return nil
}

func (r *PostRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	query := "SELECT id, title, content, author_id, author_name, image_url, created_at, expires_at FROM posts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var post domain.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.ImageURL, &post.CreatedAt, &post.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post with id %d not found", id)
		}
		return nil, fmt.Errorf("error scanning post: %w", err)
	}

	commentQuery := "SELECT id, content, author_id, image_url FROM comments WHERE post_id = $1"
	rows, err := r.db.QueryContext(ctx, commentQuery, id)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.AuthorID, &comment.ImageURL); err != nil {
			return nil, fmt.Errorf("unable to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	post.Comments = comments

	return &post, nil
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	deleteCommentsQuery := "DELETE FROM comments WHERE post_id = $1"
	_, err := r.db.ExecContext(ctx, deleteCommentsQuery, id)
	if err != nil {
		return fmt.Errorf("unable to delete comments: %w", err)
	}

	deletePostQuery := "DELETE FROM posts WHERE id = $1"
	_, err = r.db.ExecContext(ctx, deletePostQuery, id)
	if err != nil {
		return fmt.Errorf("unable to delete post: %w", err)
	}

	return nil
}

func (r *PostRepository) FindExpired(ctx context.Context) ([]*domain.Post, error) {
	query := "SELECT id, title, content, author_id, image_url FROM posts WHERE created_at < NOW() - INTERVAL '15 minutes'"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch expired posts: %w", err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.ImageURL); err != nil {
			return nil, fmt.Errorf("unable to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) FindAll(ctx context.Context) ([]*domain.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, author_id, author_name, image_url, created_at, expires_at FROM posts
	`)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch posts: %w", err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.ImageURL, &post.CreatedAt, &post.ExpiresAt); err != nil {
			return nil, fmt.Errorf("unable to scan post: %w", err)
		}

		commentRows, err := r.db.QueryContext(ctx, `
			SELECT id, post_id, title, content, author_id, author_name, image_url, reply_to_comment_id, created_at
			FROM comments WHERE post_id = $1
		`, post.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch comments for post %d: %w", post.ID, err)
		}

		var comments []domain.Comment
		for commentRows.Next() {
			var comment domain.Comment
			if err := commentRows.Scan(
				&comment.ID,
				&comment.PostID,
				&comment.Title,
				&comment.Content,
				&comment.AuthorID,
				&comment.AuthorName,
				&comment.ImageURL,
				&comment.ReplyToCommentID,
				&comment.CreatedAt,
			); err != nil {
				commentRows.Close()
				return nil, fmt.Errorf("unable to scan comment: %w", err)
			}
			comments = append(comments, comment)
		}
		commentRows.Close()

		post.Comments = comments
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) ArchiveExpiredPosts(ctx context.Context) error {
	// Implement logic for archiving expired posts
	return nil
}

// Implement the missing FindByAuthorID method
func (r *PostRepository) FindByAuthorID(ctx context.Context, authorID string) ([]*domain.Post, error) {
	query := "SELECT id, title, content, author_id, author_name, image_url, created_at, expires_at FROM posts WHERE author_id = $1"
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch posts for author %s: %w", authorID, err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.ImageURL, &post.CreatedAt, &post.ExpiresAt); err != nil {
			return nil, fmt.Errorf("unable to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) UpdateAuthorNameForPostAndComments(ctx context.Context, postID int, newAuthorName string) error {
	updatePostQuery := "UPDATE posts SET author_name = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, updatePostQuery, newAuthorName, postID)
	if err != nil {
		return fmt.Errorf("unable to update post author name: %w", err)
	}

	updateCommentsQuery := "UPDATE comments SET author_name = $1 WHERE post_id = $2"
	_, err = r.db.ExecContext(ctx, updateCommentsQuery, newAuthorName, postID)
	if err != nil {
		return fmt.Errorf("unable to update comments author name: %w", err)
	}

	return nil
}

func (r *PostRepository) UpdatePostAuthorName(ctx context.Context, postID int, authorName string) error {
	// Prepare the SQL query
	query := `UPDATE posts SET author_name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, authorName, postID)
	if err != nil {
		return fmt.Errorf("failed to update post author name: %w", err)
	}
	return nil
}

// UpdateCommentsAuthorName updates the author name of all comments for a specific post in the database
func (r *PostRepository) UpdateCommentsAuthorName(ctx context.Context, postID int, authorName string) error {
	// Prepare the SQL query
	query := `UPDATE comments SET author_name = $1 WHERE post_id = $2`
	_, err := r.db.ExecContext(ctx, query, authorName, postID)
	if err != nil {
		return fmt.Errorf("failed to update comments author name: %w", err)
	}
	return nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
        UPDATE posts
        SET author_name = $1
        WHERE id = $2
    `
	_, err := r.db.ExecContext(ctx, query, post.AuthorName, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post %d: %w", post.ID, err)
	}
	return nil
}

func (r *PostRepository) UpdateExpiration(ctx context.Context, postID int, newExpiration time.Time) error {
	query := `UPDATE posts SET expires_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newExpiration, postID)
	return err
}
