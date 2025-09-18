CREATE TABLE "customers" (
    "id" uuid  PRIMARY KEY NOT NULL  DEFAULT gen_random_uuid(),
    "name" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "name" ON "customers" ("name");