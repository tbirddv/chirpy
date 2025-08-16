-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW() WHERE id = $3 RETURNING *;

-- name: GiveChirpyRed :one
UPDATE users SET is_chirpy_red = true, updated_at = NOW() WHERE id = $1 RETURNING *;
