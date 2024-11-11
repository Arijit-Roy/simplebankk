-- name: CreateUser :one
INSERT INTO users (
  user_name, hashed_password, full_name, email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE user_name=$1  LIMIT 1;

-- -- name: GetAccountFroUpdate :one
-- SELECT * FROM accounts 
-- WHERE id=$1  LIMIT 1 FOR NO KEY UPDATE;
