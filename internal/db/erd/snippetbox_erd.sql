CREATE TABLE "snippets" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "content" text NOT NULL,
  "created" timestamptz NOT NULL DEFAULT (now()),
  "expires" timestamptz NOT NULL
);

CREATE INDEX "idx_snippets_created" ON "snippets" ("created");
