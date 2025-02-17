-- name: CreateBookReview :one
INSERT INTO book_reviews (
    id, book_id, rating, text
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;
