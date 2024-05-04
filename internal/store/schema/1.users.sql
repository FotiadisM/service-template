CREATE TYPE user_scope AS ENUM ('user', 'admin');

CREATE TABLE users (
    id UUID NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    scope USER_SCOPE NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users (email);
