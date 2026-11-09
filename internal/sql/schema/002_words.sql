-- +goose Up
create table words(
    id serial primary key,
    word text unique
);

-- +goose Down
drop table words;
