-- name: GetPostsforUser :many
SELECT a.*, b.name from posts as a
INNER JOIN feeds as b on a.feed_id = b.id
INNER JOIN feed_follows as c on b.id = c.feed_id
WHERE c.user_id = $1
ORDER BY published_at DESC
LIMIT $2;
