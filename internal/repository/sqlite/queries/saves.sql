-- name: ListSaves :many
SELECT id, user_id, name, saved_at FROM saves WHERE user_id = ? ORDER BY saved_at DESC;

-- name: GetSave :one
SELECT * FROM saves WHERE id = ? AND user_id = ?;

-- name: CreateSave :one
INSERT INTO saves (user_id, name, state) VALUES (?, ?, ?) RETURNING *;

-- name: UpdateSave :exec
UPDATE saves SET name = ?, state = ?, saved_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?;

-- name: DeleteSave :exec
DELETE FROM saves WHERE id = ? AND user_id = ?;
