// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package queries

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID
	Name      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Book struct {
	ID          uuid.UUID
	Title       string
	AuthorID    uuid.UUID
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
