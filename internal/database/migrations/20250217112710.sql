-- Create "book_reviews" table
CREATE TABLE "public"."book_reviews" (
    "id" uuid NOT NULL,
    "book_id" uuid NOT NULL,
    "rating" integer NOT NULL,
    "text" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    CONSTRAINT "book_reviews_book_id_fkey" FOREIGN KEY (
        "book_id"
    ) REFERENCES "public"."books" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
