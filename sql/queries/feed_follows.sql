-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (user_id, feed_id)
    VALUES ($1, $2)
    RETURNING *
)
SELECT
    inserted_feed_follow.id,
    inserted_feed_follow.user_id,
    inserted_feed_follow.feed_id,
    inserted_feed_follow.created_at,
    inserted_feed_follow.updated_at,
    users.name AS user_name,
    feed.name AS feed_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feed ON feed.id = inserted_feed_follow.feed_id;


-- name: GetFeedFollowsForUser :many
WITH selected_user AS (
    SELECT id
    FROM users
    WHERE users.name = $1
)
SELECT 
    feed_follows.id,
    feed_follows.user_id,
    feed_follows.feed_id,
    feed_follows.created_at,
    feed_follows.updated_at,
    users.name AS user_name,   
    feed.name AS feed_name    
FROM feed_follows
JOIN selected_user ON feed_follows.user_id = selected_user.id
JOIN users ON feed_follows.user_id = users.id 
JOIN feed ON feed_follows.feed_id = feed.id;
