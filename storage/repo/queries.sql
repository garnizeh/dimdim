-- name: ListAllTags :many
SELECT *
FROM tags
ORDER BY name;

-- name: GetTagByID :one
SELECT *
FROM tags
WHERE id = ?1 AND deleted_at = 0
LIMIT 1;

-- name: GetTagByName :one
SELECT *
FROM tags
WHERE name = ?1 AND deleted_at = 0
LIMIT 1;

-- name: CreateTag :exec
INSERT INTO tags (
  name, created_at, updated_at
) VALUES (
  ?1, ?2, ?2
);

-- name: UpdateTag :exec
UPDATE tags
set name = ?2, updated_at = ?3
WHERE id = ?1 AND deleted_at = 0;

-- name: DeleteTag :exec
UPDATE tags
set deleted_at = ?2
WHERE id = ?1 AND deleted_at = 0;