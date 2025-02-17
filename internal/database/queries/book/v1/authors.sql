-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors;

-- name: CreateAuthor :one
INSERT INTO authors (
    id, name, bio
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- name: UpdateAuthor :one
UPDATE authors
SET
    name = coalesce(sqlc.narg('name'), name),
    bio = coalesce(sqlc.narg('bio'), bio),
    updated_at = sqlc.arg('updated_at')
WHERE id = $1
RETURNING *;
