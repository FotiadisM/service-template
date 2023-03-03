CREATE TYPE user_scope AS ENUM (
    'applicant', 'company'
);

CREATE TABLE users (
    id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT id PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users(user_id);

CREATE TABLE user_scope (
    user_id UUID NOT NULL,
    scope USER_SCOPE NOT NULL,

    CONSTRAINT id PRIMARY KEY (user_id)
);

ALTER TABLE user_scope
ADD CONSTRAINT user_scope_user_id_fkey FOREIGN KEY (
    user_id
) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;
