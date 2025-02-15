CREATE TABLE authors (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    bio TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT authors_pkey PRIMARY KEY (id)
);
