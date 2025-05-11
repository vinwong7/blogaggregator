-- name: Unfollow :exec
DELETE FROM feed_follows
where feed_id = $1
and user_id = $2;