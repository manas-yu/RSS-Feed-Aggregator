-- +goose Up
ALTER TABLE users ADD column api_keys VARCHAR(64) UNIQUE NOT NULL DEFAULT(
    encode(sha256(random()::text::bytea), 'hex')
);
-- +goose Down
ALTER TABLE users DROP column api_keys;
