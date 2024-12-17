CREATE TABLE post_comments (
    id UUID NOT NULL,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    likes INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT post_comments_pkey PRIMARY KEY (id),
    CONSTRAINT post_comments_post_id_fkey FOREIGN KEY (
        post_id
    ) REFERENCES posts (id),
    CONSTRAINT post_comments_user_id_fkey FOREIGN KEY (
        user_id
    ) REFERENCES users (id)
);
