CREATE TYPE user_scope AS ENUM (
    'applicant', 'company', 'admin'
);

CREATE TABLE users (
    id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    scope USER_SCOPE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT id PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users(email);
