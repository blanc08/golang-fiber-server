CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "email" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expired_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT 'now()'
);
-- Add Foreign key
ALTER TABLE "sessions"
ADD FOREIGN KEY ("email") REFERENCES "users" ("email");