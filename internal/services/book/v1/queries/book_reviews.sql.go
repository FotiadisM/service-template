// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: book_reviews.sql

package queries

import (
	"context"

	"github.com/google/uuid"
)

const createBookReview = `-- name: CreateBookReview :one
INSERT INTO book_reviews (
    id, book_id, rating, text
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, book_id, rating, text, created_at, updated_at
`

type CreateBookReviewParams struct {
	ID     uuid.UUID
	BookID uuid.UUID
	Rating int32
	Text   string
}

func (q *Queries) CreateBookReview(ctx context.Context, arg CreateBookReviewParams) (BookReview, error) {
	row := q.db.QueryRowContext(ctx, createBookReview,
		arg.ID,
		arg.BookID,
		arg.Rating,
		arg.Text,
	)
	var i BookReview
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.Rating,
		&i.Text,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
