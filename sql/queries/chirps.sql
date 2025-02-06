/*
 * Copyright 2025 Merlinux-source
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
