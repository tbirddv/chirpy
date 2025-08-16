-- name: CreateChirp :one
INSERT into chirps (id, user_id, body) 
values (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * from chirps
ORDER BY created_at ASC;

-- name: GetChirpsDesc :many
SELECT * from chirps
ORDER BY created_at DESC;

-- name: GetChirpsByUser :many
SELECT * from chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsByUserDesc :many
SELECT * from chirps
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetChirpByID :one
SELECT * from chirps where id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;