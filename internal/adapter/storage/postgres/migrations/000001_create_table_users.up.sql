CREATE TYPE "users_role_enum" AS ENUM ('admin', 'member');

CREATE TABLE "users" (
    "id" uuid  PRIMARY KEY NOT NULL  DEFAULT gen_random_uuid(),
    "fullname" varchar NOT NULL,
    "email" varchar NOT NULL,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL,
    "role" users_role_enum DEFAULT 'member',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "email" ON "users" ("email");

INSERT INTO users (id, fullname, email, username, password, role, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  'Administrador',
  'admin@admin.com', 
  'admin', 
  '$2a$10$wvyY/NTJ4PYnxpx8MhrGO.wWHRjwKAbNUUbSXTMkzeWxBMx8oS9K.', 
  'admin', 
  current_timestamp,
  current_timestamp
  );
