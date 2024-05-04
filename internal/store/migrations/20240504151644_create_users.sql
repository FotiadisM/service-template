-- Create enum type "user_scope"
CREATE TYPE "public"."user_scope" AS ENUM ('user', 'admin');

-- Create "users" table
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL,
    "email" text NOT NULL,
    "password" text NOT NULL,
    "scope" "public"."user_scope" NOT NULL,
    "active" boolean NOT NULL DEFAULT true,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY ("id")
);

-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
