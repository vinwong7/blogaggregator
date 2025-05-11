-- name: GetFeedFollowsForUser :many
SELECT b.name as feed_name, c.name as user_name from feed_follows as a
INNER JOIN feeds as b on a.feed_id = b.id
INNER JOIN users as c on a.user_id = c.id
where c.name = $1;

