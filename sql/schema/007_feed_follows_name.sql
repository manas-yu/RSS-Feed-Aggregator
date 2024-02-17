-- +goose Up
ALTER TABLE feed_follows
ADD column name VARCHAR(64);
-- +goose Down
ALTER TABLE feed_follows DROP column name;