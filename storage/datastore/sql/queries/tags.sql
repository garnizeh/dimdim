-- name: CreateTag :exec
INSERT INTO tags (name)
          VALUES (?);

-- name: DeleteTag :exec
UPDATE tags SET updated_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER), deleted_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE id = ?1 AND deleted_at = 0;

-- name: UpdateTag :exec
UPDATE tags
SET name = ?, updated_at = CAST(unixepoch('subsecond') * 1000 as int)
WHERE id = ?;

-- name: ListAllTags :many
SELECT * FROM tags
WHERE deleted_at = 0
ORDER BY name;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = ? AND deleted_at = 0;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE name = ? AND deleted_at = 0;
