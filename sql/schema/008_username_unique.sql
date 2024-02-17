-- +goose Up
ALTER TABLE users
ADD CONSTRAINT unique_name UNIQUE (name);
-- +goose Down
-- Down migration is not needed for adding a unique constraint