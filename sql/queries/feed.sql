-- name: CreateFeed :one
INSERT INTO feed (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListAllFeeds :many
SELECT 
    feed.id,
    feed.name AS feed_name,
    feed.url,
    feed.user_id,
    users.name AS user_name,
    feed.created_at,
    feed.updated_at
FROM feed
JOIN users ON feed.user_id = users.id;
