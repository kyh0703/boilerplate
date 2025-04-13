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
  update_at,
  create_at
) VALUES (
  ?, ?, ?, ?, now(), now()
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET
email = ?,
name = ?,
password = ?,
bio = ?,
update_at = now()
WHERE id = ?
RETURNING *;

-- name: PatchUser :exec
UPDATE users SET
name = COALESCE(sqlc.narg(name), name),
password = COALESCE(sqlc.narg(password), password),
bio = COALESCE(sqlc.narg(bio), bio),
update_at = now()
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
  ?, ?, ?, now()
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

-- name: CreateProject :one
INSERT INTO projects (
  user_id,
  name,
  description,
  update_at,
  create_at
) VALUES (
  ?, ?, ?, now(), now()
)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = ? LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
WHERE user_id = ?
ORDER BY name;

-- name: PatchProject :exec
UPDATE projects SET
name = COALESCE(sqlc.narg(name), name),
description = COALESCE(sqlc.narg(description), description),
update_at = now()
WHERE id = ?
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ?;

-- name: CreateFlow :one
INSERT INTO flows (
  name,
  description
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetFlow :one
SELECT * FROM flows
WHERE id = ? LIMIT 1;

-- name: ListFlows :many
SELECT * FROM flows
WHERE project_id = ?
ORDER BY name;

-- name: PatchFlow :exec
UPDATE flows SET
name = COALESCE(sqlc.narg(name), name),
description = COALESCE(sqlc.narg(description), description),
update_at = now()
WHERE id = ?
RETURNING *;

-- name: DeleteFlow :exec
DELETE FROM flows
WHERE id = ?;

-- name: GetNode :one
SELECT * FROM nodes
WHERE id = ? LIMIT 1;

-- name: ListNodes :many
SELECT * FROM nodes
WHERE flow_id = ?
ORDER BY create_at;

-- name: CreateNode :one
INSERT INTO nodes (
  id,
  flow_id,
  type,
  position,
  styles,
  width,
  height,
  hidden,
  description,
  update_at,
  create_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now()
)
RETURNING *;

-- name: PatchNode :exec
UPDATE nodes SET
type = COALESCE(sqlc.narg(type), type),
position = COALESCE(sqlc.narg(position), position),
styles = COALESCE(sqlc.narg(styles), styles),
width = COALESCE(sqlc.narg(width), width),
height = COALESCE(sqlc.narg(height), height),
hidden = COALESCE(sqlc.narg(hidden), hidden),
description = COALESCE(sqlc.narg(description), description),
update_at = now()
WHERE id = ?
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes
WHERE id = ?;

-- name: GetEdge :one
SELECT * FROM edges
WHERE id = ? LIMIT 1;

-- name: ListEdges :many
SELECT * FROM edges
WHERE flow_id = ?
ORDER BY create_at;

-- name: CreateEdge :one
INSERT INTO edges (
  id,
  flow_id,
  source,
  target,
  type,
  label,
  hidden,
  marker_end,
  update_at,
  create_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, now(), now()
)
RETURNING *;

-- name: PatchEdge :exec
UPDATE edges SET
source = COALESCE(sqlc.narg(source), source),
target = COALESCE(sqlc.narg(target), target),
type = COALESCE(sqlc.narg(type), type),
label = COALESCE(sqlc.narg(label), label),
hidden = COALESCE(sqlc.narg(hidden), hidden),
marker_end = COALESCE(sqlc.narg(marker_end), marker_end),
update_at = now()
WHERE id = ?
RETURNING *;

-- name: DeleteEdge :exec
DELETE FROM edges
WHERE id = ?;
