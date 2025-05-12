-- +goose Up
ALTER TABLE feeds 
ADD last_fetched_dt TIMESTAMP;

-- +goose Down
ALTER TABLE feeds 
DROP COLUMN last_fetched_dt;
