-- Modify "authors" table
ALTER TABLE "public"."authors" ADD COLUMN "bio" text NOT NULL DEFAULT '';
-- Modify "books" table
ALTER TABLE "public"."books" ADD COLUMN "description" text NOT NULL DEFAULT '';
