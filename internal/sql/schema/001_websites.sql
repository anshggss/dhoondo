-- +goose Up
create table websites (
    id uuid primary key,
    url text unique
);

-- +goose Down
drop table websites;
