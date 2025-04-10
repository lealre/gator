-- name: CreatePost :one
INSERT INTO posts (
    id,
    title,
    url,
    description,
    feed_id,
    published_at,
    updated_at,
    created_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;


-- name: ListPostsByUser :many
SELECT 
    posts.id,
    posts.title,
    posts.url,
    posts.description,
    posts.feed_id,
    posts.published_at,
    posts.updated_at,
    posts.created_at
FROM posts
JOIN feed ON posts.feed_id = feed.id
WHERE feed.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;

