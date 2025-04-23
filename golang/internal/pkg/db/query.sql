-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  email,
  password,
  name,
  bio,
  provider,
  provider_id,
  is_admin,
  update_at,
  create_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET
email = ?,
name = ?,
password = ?,
bio = ?,
provider = ?,
provider_id = ?,
is_admin = ?,
update_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: PatchUser :exec
UPDATE users SET
name = COALESCE(sqlc.narg(name), name),
password = COALESCE(sqlc.narg(password), password),
bio = COALESCE(sqlc.narg(bio), bio),
provider = COALESCE(sqlc.narg(provider), provider),
provider_id = COALESCE(sqlc.narg(provider_id), provider_id),
is_admin = COALESCE(sqlc.narg(is_admin), is_admin),
update_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetToken :one
SELECT * FROM tokens
WHERE id = ? LIMIT 1;

-- name: GetTokenByUserID :one
SELECT * FROM tokens
WHERE user_id = ? LIMIT 1;

-- name: ListTokens :many
SELECT * FROM tokens
ORDER BY create_at;

-- name: CreateToken :one
INSERT INTO tokens (
  user_id,
  refresh_token,
  expires_in,
  create_at
) VALUES (
  ?, ?, ?, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: UpdateToken :exec
UPDATE tokens SET
user_id = ?,
refresh_token = ?,
expires_in = ?
WHERE id = ?
RETURNING *;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = ?;

-- name: CreatePost :one
INSERT INTO posts (
  user_id,
  title,
  content,
  update_at,
  create_at
) VALUES (
  ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ? LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
WHERE user_id = ?
ORDER BY title;

-- name: ListPostsWithPaging :many
SELECT * FROM posts
WHERE user_id = ?
ORDER BY title
LIMIT ? OFFSET ?;

-- name: PatchPost :exec
UPDATE posts SET
title = COALESCE(sqlc.narg(title), title),
content = COALESCE(sqlc.narg(content), content),
update_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = ?;

-- name: CreateOAuthState :one
INSERT INTO oauth_states (
  state,
  redirect_url,
  expires_at,
  create_at
) VALUES (
  ?, ?, ?, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetOAuthState :one
SELECT * FROM oauth_states
WHERE state = ? LIMIT 1;

-- name: DeleteOAuthState :exec
DELETE FROM oauth_states
WHERE state = ?;
