-- Create "books" table
CREATE TABLE "public"."books" (
    "id" uuid NOT NULL,
    "title" text NOT NULL,
    "author_id" uuid NOT NULL,
    "description" text NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    CONSTRAINT "books_author_id_fkey" FOREIGN KEY (
        "author_id"
    ) REFERENCES "public"."authors" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
