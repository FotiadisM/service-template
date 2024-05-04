CREATE TABLE posts (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    likes INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT posts_pkey PRIMARY KEY (id),
    CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id)
);
