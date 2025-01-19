-- Create "post_comments" table
CREATE TABLE "public"."post_comments" (
    "id" uuid NOT NULL,
    "post_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "text" text NOT NULL,
    "likes" integer NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY ("id"),
    CONSTRAINT "post_comments_post_id_fkey" FOREIGN KEY (
        "post_id"
    ) REFERENCES "public"."posts" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT "post_comments_user_id_fkey" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
