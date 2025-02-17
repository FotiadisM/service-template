CREATE TABLE book_reviews (
    id UUID NOT NULL,
    book_id UUID NOT NULL,
    rating INTEGER NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT book_reviews_pkey PRIMARY KEY (id),
    CONSTRAINT book_reviews_book_id_fkey FOREIGN KEY (
        book_id
    ) REFERENCES books (id) ON UPDATE CASCADE ON DELETE CASCADE
);
