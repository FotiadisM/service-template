-- Create "authors" table
CREATE TABLE "public"."authors" (
    "id" uuid NOT NULL,
    "name" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
