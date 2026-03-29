CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "pos_id" int NOT NULL,
  "price" int NOT NULL,
  "recipe_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
