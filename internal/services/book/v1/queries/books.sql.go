// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: books.sql

package queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createBook = `-- name: CreateBook :one
INSERT INTO books (
    id, title, author_id, description
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, title, author_id, description, created_at, updated_at
`

type CreateBookParams struct {
	ID          uuid.UUID
	Title       string
	AuthorID    uuid.UUID
	Description string
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, createBook,
		arg.ID,
		arg.Title,
		arg.AuthorID,
		arg.Description,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.AuthorID,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteBook = `-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1
`

func (q *Queries) DeleteBook(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteBook, id)
	return err
}

const getBook = `-- name: GetBook :one
SELECT id, title, author_id, description, created_at, updated_at FROM books
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBook(ctx context.Context, id uuid.UUID) (Book, error) {
	row := q.db.QueryRowContext(ctx, getBook, id)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.AuthorID,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listBooks = `-- name: ListBooks :many
SELECT id, title, author_id, description, created_at, updated_at FROM books
`

func (q *Queries) ListBooks(ctx context.Context) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, listBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Book{}
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.AuthorID,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBook = `-- name: UpdateBook :one
UPDATE books
SET
    title = coalesce($2, title),
    description = coalesce($3, description),
    updated_at = $4
WHERE id = $1
RETURNING id, title, author_id, description, created_at, updated_at
`

type UpdateBookParams struct {
	ID          uuid.UUID
	Title       sql.NullString
	Description sql.NullString
	UpdatedAt   time.Time
}

func (q *Queries) UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, updateBook,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.UpdatedAt,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.AuthorID,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
