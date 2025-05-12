-- name: MarkFeedFetched :exec
UPDATE feeds
set last_fetched_dt = cast(Now() as TIMESTAMP),
updated_at = cast(Now() as TIMESTAMP)
where id = $1;