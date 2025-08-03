-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-08-03T10:26:16.931Z

CREATE TYPE "role" AS ENUM (
  'owner',
  'editor',
  'viewer'
);

CREATE TABLE "Book" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Category" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Redial" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "date" date NOT NULL,
  "time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Customer" (
  "id" uuid PRIMARY KEY,
  "book_id" uuid NOT NULL,
  "category_id" uuid,
  "contact_id" uuid UNIQUE NOT NULL,
  "redial_id" uuid UNIQUE,
  "name" varchar NOT NULL,
  "corporation" varchar,
  "address" varchar,
  "leader" uuid UNIQUE,
  "pic" uuid UNIQUE,
  "memo" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Staff" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "sex" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Contact" (
  "id" uuid PRIMARY KEY,
  "customer_id" uuid NOT NULL,
  "staff_id" uuid UNIQUE,
  "phone" varchar NOT NULL,
  "mail" varchar,
  "fax" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Call" (
  "id" uuid PRIMARY KEY,
  "customer_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "status_id" uuid,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Status" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "effective" bool,
  "ng" bool,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "User" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "name" varchar NOT NULL,
  "role" role NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "Customer"."leader" IS '代表者';

COMMENT ON COLUMN "Customer"."pic" IS '担当者';

COMMENT ON COLUMN "Contact"."staff_id" IS 'これがnullの場合代表';

COMMENT ON COLUMN "Status"."effective" IS '有効数としてカウントするか';

COMMENT ON COLUMN "Status"."ng" IS 'NG';

ALTER TABLE "Customer" ADD FOREIGN KEY ("book_id") REFERENCES "Book" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "Customer" ADD FOREIGN KEY ("category_id") REFERENCES "Category" ("id");

ALTER TABLE "Call" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id");

ALTER TABLE "Call" ADD FOREIGN KEY ("status_id") REFERENCES "Status" ("id");

ALTER TABLE "Call" ADD FOREIGN KEY ("customer_id") REFERENCES "Customer" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "Customer" ADD FOREIGN KEY ("leader") REFERENCES "Staff" ("id");

ALTER TABLE "Customer" ADD FOREIGN KEY ("pic") REFERENCES "Staff" ("id");

ALTER TABLE "Staff" ADD FOREIGN KEY ("id") REFERENCES "Contact" ("staff_id");

ALTER TABLE "Contact" ADD FOREIGN KEY ("customer_id") REFERENCES "Customer" ("id");

ALTER TABLE "Contact" ADD FOREIGN KEY ("id") REFERENCES "Customer" ("contact_id");

ALTER TABLE "Redial" ADD FOREIGN KEY ("id") REFERENCES "Customer" ("redial_id");

ALTER TABLE "Redial" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id");
