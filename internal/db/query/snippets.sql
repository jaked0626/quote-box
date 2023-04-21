-- name: CreateSnippet :one
INSERT INTO snippets (
    title,
    content,
    expires
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetSnippetByID :one 
SELECT * FROM snippets
WHERE id = $1 
LIMIT 1;

-- name: ListLatestSnippets :many
SELECT * FROM snippets
ORDER BY created DESC
LIMIT $1;