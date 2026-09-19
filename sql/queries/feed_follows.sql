



-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feeds_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds_follows.id, feeds_follows.created_at, feeds_follows.updated_at, feeds_follows.user_id, feeds_follows.feed_id,
    feeds.name as feed_name,
    users.name as user_name
FROM feeds_follows
INNER JOIN feeds ON feeds_follows.feed_id = feeds.id
INNER JOIN users ON feeds_follows.user_id = users.id
WHERE feeds_follows.user_id = $1;