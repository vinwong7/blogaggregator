-- name: FeedList :many
SELECT a.name as feedName, a.url, b.name as userName 
FROM feeds as a 
INNER JOIN users as b on a.user_id = b.id;
