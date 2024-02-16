-- +goose Up
create table users(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar(255) not null --TODO: add unique constraint
);
-- +goose Down
drop table users;