-- name: GetNextFeedtoFetch :one
SELECT * FROM feeds 
ORDER BY last_fetched_dt ASC NULLS FIRST
LIMIT 1
;
