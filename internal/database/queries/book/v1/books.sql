-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books;

-- name: CreateBook :one
INSERT INTO books (
    id, title, author_id, description
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: UpdateBook :one
UPDATE books
SET
    title = coalesce(sqlc.narg('title'), title),
    description = coalesce(sqlc.narg('description'), description),
    updated_at = sqlc.arg('updated_at')
WHERE id = $1
RETURNING *;
