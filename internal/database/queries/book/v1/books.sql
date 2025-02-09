-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books;

-- name: CreateBook :exec
INSERT INTO books (
    id, title, author_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;
