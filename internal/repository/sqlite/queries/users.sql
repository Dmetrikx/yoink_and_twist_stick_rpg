-- name: UpsertUser :exec
INSERT INTO users (id, email) VALUES (?, ?)
ON CONFLICT(id) DO UPDATE SET email=excluded.email;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
