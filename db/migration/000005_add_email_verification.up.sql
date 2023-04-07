ALTER TABLE users ADD is_email_verified boolean NOT NULL DEFAULT FALSE;

CREATE TABLE "verify_emails"(
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL REFERENCES "users" ("username"),
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" boolean NOT NULL DEFAULT FALSE,
    "create_at" timestamptz NOT NULL DEFAULT (now()),
    "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);