-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors;

-- name: CreateAuthor :exec
INSERT INTO authors (
  id, name, created_at, updated_at
) VALUES (
  $1, $2, $3, $4
);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
