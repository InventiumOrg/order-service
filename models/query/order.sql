-- name: CreateOrder :one
INSERT INTO orders (
    pos_id, price, recipe_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET pos_id = $2,
    price = $3,
    recipe_id = $4
WHERE id = $1
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1;

-- name: ListOrder :many
SELECT id, pos_id, price, recipe_id
FROM orders
LIMIT $1 OFFSET $2;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
