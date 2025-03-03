-- name: CreateUser :exec
INSERT INTO users (id, email, name, password, salt)
           VALUES (? , ?    , ?   , ?       , ?);

-- name: DeleteUser :exec
UPDATE users SET updated_at = CAST(unixepoch('subsecond') * 1000 as int), deleted_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users SET email = ?, name = ?, updated_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users SET password = ?, salt = ?, updated_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE id = ?;

-- name: SetUserIsVerified :one
UPDATE users SET verified_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER), updated_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE email = ?
RETURNING *;

-- name: GetUserByID :one
SELECT  * FROM users
WHERE id = ? AND deleted_at = 0;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? AND deleted_at = 0;

-- name: GetUserIsVerified :one
SELECT true FROM users
WHERE id = ? AND deleted_at = 0 AND verified_at > 0;

-- name: GetUserIsDeleted :one
SELECT true FROM users
WHERE id = ? AND deleted_at > 0;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE id = ? AND deleted_at > 0
ORDER BY name;