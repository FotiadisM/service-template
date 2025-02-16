package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/FotiadisM/service-template/internal/services/book/v1/queries"
)

func uuidParser(t *testing.T, s string) uuid.UUID {
	t.Helper()
	u, err := uuid.Parse(s)
	require.NoError(t, err, "failed to parse uuid")
	return u
}

type Fixtures struct {
	Author1 queries.Author
	Book1   queries.Book
	Book2   queries.Book
	Book3   queries.Book

	Author2 queries.Author
	Book4   queries.Book
	Book5   queries.Book
}

func NewFixtures(t *testing.T) *Fixtures {
	author1 := queries.Author{
		ID:        uuidParser(t, "01950b20-756a-7056-bd6a-272db29cb3d1"),
		Name:      "Author1",
		Bio:       "This is author's 1 Bio",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	author2 := queries.Author{
		ID:        uuidParser(t, "01950b20-756a-7ed2-a49b-adfc1d46533b"),
		Name:      "Author2",
		Bio:       "This is author's 2 Bio",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &Fixtures{
		Author1: author1,
		Book1: queries.Book{
			ID:          uuidParser(t, "01950b20-756a-730a-8816-a7d8a675fc3e"),
			Title:       "Book1",
			AuthorID:    author1.ID,
			Description: "This is book 1 description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Book2: queries.Book{
			ID:          uuidParser(t, "01950b20-756a-7c88-9bf5-91c8e6216938"),
			Title:       "Book2",
			AuthorID:    author1.ID,
			Description: "This is book 2 description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Book3: queries.Book{
			ID:          uuidParser(t, "01950b20-756a-7a24-8fa3-b25971286e08"),
			Title:       "Book3",
			AuthorID:    author1.ID,
			Description: "This is book 3 description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Author2: author2,
		Book4: queries.Book{
			ID:          uuidParser(t, "01950b20-756a-726d-9a9b-0d0bbeba6b31"),
			Title:       "Book4",
			AuthorID:    author2.ID,
			Description: "This is book 4 description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Book5: queries.Book{
			ID:          uuidParser(t, "01950b20-756a-7f3d-b2b0-9bb330e015c3"),
			Title:       "Book5",
			AuthorID:    author2.ID,
			Description: "This is book 5 description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

func authorToCreateParams(author queries.Author) queries.CreateAuthorParams {
	return queries.CreateAuthorParams{
		ID:        author.ID,
		Name:      author.Name,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}
}

func bookToCreateParams(book queries.Book) queries.CreateBookParams {
	return queries.CreateBookParams{
		ID:          book.ID,
		Title:       book.Title,
		AuthorID:    book.AuthorID,
		Description: book.Description,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}

func (f *Fixtures) Load(ctx context.Context, t *testing.T, db *sql.DB) {
	authors := []queries.Author{f.Author1, f.Author2}
	books := []queries.Book{f.Book1, f.Book2, f.Book3, f.Book4, f.Book5}

	querier := queries.New(db)

	var err error
	for _, author := range authors {
		_, err = querier.CreateAuthor(ctx, authorToCreateParams(author))
		require.NoError(t, err, "failed to create author")
	}

	for _, book := range books {
		_, err = querier.CreateBook(ctx, bookToCreateParams(book))
		require.NoError(t, err, "failed to create book")
	}
}
