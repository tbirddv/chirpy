-- name: CreateChirp :one
INSERT into chirps (id, user_id, body) 
values (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * from chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * from chirps where id = $1;