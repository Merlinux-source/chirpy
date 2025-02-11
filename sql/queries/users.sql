-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES ($1, $2)
RETURNING *;

-- name: ClearUsers :exec
DELETE FROM USERS;
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: SetUserEmailAndPassword :one
UPDATE users
SET email           = $2,
    hashed_password = $3
WHERE id = $1
RETURNING *;