CREATE TABLE books (
    id UUID NOT NULL,
    title TEXT NOT NULL,
    author_id UUID NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT books_pkey PRIMARY KEY (id),
    CONSTRAINT books_author_id_fkey FOREIGN KEY (
        author_id
    ) REFERENCES authors (id) ON UPDATE CASCADE ON DELETE CASCADE
);
