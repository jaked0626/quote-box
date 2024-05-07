CREATE TABLE "snippets" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "title" VARCHAR NOT NULL,
  "author" VARCHAR NOT NULL DEFAULT 'Unknown',
  "work" VARCHAR NOT NULL DEFAULT 'Unknown',
  "content" TEXT NOT NULL,
  "created" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "expires" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '100 days')
);

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "hashed_password" CHAR(60) NOT NULL,
  "created" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "sessions" (
  "token" TEXT PRIMARY KEY NOT NULL,
  "data" BYTEA NOT NULL,
  "expiry" TIMESTAMPTZ NOT NULL
);

CREATE INDEX "idx_snippets_created" ON "snippets" ("created");

CREATE INDEX "sessions_expiry_idx" ON "sessions" ("expiry");