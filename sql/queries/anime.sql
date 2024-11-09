-- name: GetAnime :one
SELECT * FROM anime
WHERE id = ? LIMIT 1;

-- name: GetAllAnime :many
SELECT * FROM anime;

-- name: CreateAnime :one
INSERT INTO anime (id, title, startDate, updatedDate, completion)
VALUES (?, ?, ?, ?, ?)
RETURNING *;
