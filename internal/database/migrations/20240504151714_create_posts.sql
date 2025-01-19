-- Create "posts" table
CREATE TABLE "public"."posts" (
    "id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "title" text NOT NULL,
    "text" text NOT NULL,
    "likes" integer NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY ("id"),
    CONSTRAINT "posts_user_id_fkey" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
