-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2) RETURNING *;
-- name: ClearChirps :exec
DELETE
FROM chirps;
-- name: GetChrips :many
SELECT *
FROM chirps
ORDER BY created_at ASC;
-- name: GetChripsFromTo :many
SELECT *
FROM chirps
WHERE updated_at BETWEEN $1 AND $2;
-- name: GetChripsByUserIdFromTo :many
SELECT *
FROM chirps
WHERE user_id = $1
  AND updated_at BETWEEN $2 AND $3;
-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;
-- name: GetChirpsByUserId :many
SELECT *
FROM chirps
WHERE user_id = $1;
-- name: DeleteChirp :exec
DELETE
FROM chirps
WHERE id = $1;
