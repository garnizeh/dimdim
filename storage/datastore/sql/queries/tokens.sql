-- name: CreateToken :exec
INSERT INTO tokens (token, type, email, expires_at)
            VALUES (?    , ?   , ?    , ?);

-- name: DeleteExpiredTokens :exec
UPDATE tokens SET deleted_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE expires_at <= ?;

-- name: DeleteSignupTokensByEmail :exec
UPDATE tokens SET deleted_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE email = ? AND type = 'SIGNUP';

-- name: DeletePasswordTokensByEmail :exec
UPDATE tokens SET deleted_at = CAST(unixepoch('subsecond') * 1000 AS INTEGER)
WHERE email = ? AND type = 'PASSWORD';

-- name: GetSignupTokenNotExpired :one
SELECT * FROM tokens
WHERE token = ? AND type = 'SIGNUP' AND expires_at >= ? AND deleted_at = 0;

-- name: GetPasswordTokenNotExpired :one
SELECT * FROM tokens
WHERE token = ? AND type = 'PASSWORD' AND expires_at >= ? AND deleted_at = 0;