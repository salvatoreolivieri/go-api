package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	UserId    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {

	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM posts
		WHERE id = $1
	`
	var post Post

	err := s.db.QueryRowContext(
		ctx,
		query,
		id).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound

		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) DeleteByID(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts
		where id = $1
	`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Optionally check if any rows were affected
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted with id: %d", id)
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {

	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, content, updated_at
	`

	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	return nil
}
