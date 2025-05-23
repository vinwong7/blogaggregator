// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: markFeedFetched.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
set last_fetched_dt = cast(Now() as TIMESTAMP),
updated_at = cast(Now() as TIMESTAMP)
where id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
